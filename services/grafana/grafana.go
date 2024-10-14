// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// Package grafana provides facilities for working with Grafana.
package grafana

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
)

const (
	defaultOrgID              = 1
	defaultAlertRuleNamespace = "Insight"
	defaultDatasource         = "Prometheus"
	defaultIntervalSeconds    = 60
)

// Client represents a client for Grafana API.
type Client struct {
	addr       string
	db         string
	alertsPath string
	http       *http.Client
}

// NewClient creates a new client for given Grafana address.
func NewClient(addr string, db string, alertsPath string) *Client {
	return &Client{
		addr:       addr,
		db:         db,
		alertsPath: alertsPath,
		http:       &http.Client{},
	}
}

type annotation struct {
	Time time.Time `json:"-"`
	Tags []string  `json:"tags,omitempty"`
	Text string    `json:"text,omitempty"`

	TimeInt int64 `json:"time,omitempty"`
}

// encode annotation before sending request.
func (a *annotation) encode() {
	var t int64
	if !a.Time.IsZero() {
		t = a.Time.UnixNano() / int64(time.Millisecond)
	}
	a.TimeInt = t
}

// decode annotation after receiving response.
func (a *annotation) decode() {
	var t time.Time
	if a.TimeInt != 0 {
		t = time.Unix(0, a.TimeInt*int64(time.Millisecond))
	}
	a.Time = t
}

// CreateAnnotation creates annotation with given text and tags ("pmm_annotation" is added automatically)
// and returns Grafana's response text which is typically "Annotation added" or "Failed to save annotation".
func (c *Client) CreateAnnotation(ctx context.Context, tags []string, text string) (string, error) {
	// http://docs.grafana.org/http_api/annotations/#create-annotation

	request := &annotation{
		Tags: append([]string{"pmm_annotation"}, tags...),
		Text: text,
	}
	request.encode()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return "", errors.Wrap(err, "failed to marhal request")
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.addr,
		Path:   "/api/annotations",
	}
	resp, err := c.http.Post(u.String(), "application/json", &buf)
	if err != nil {
		return "", errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Get(ctx).Warnf("Grafana responded with status %d.", resp.StatusCode)
	}

	var response struct {
		Message string `json:"message"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", errors.Wrap(err, "failed to decode JSON response")
	}
	return response.Message, nil
}

func (c *Client) findAnnotations(ctx context.Context, from, to time.Time) ([]annotation, error) {
	// http://docs.grafana.org/http_api/annotations/#find-annotations

	u := &url.URL{
		Scheme: "http",
		Host:   c.addr,
		Path:   "/api/annotations",
		RawQuery: url.Values{
			"from": []string{strconv.FormatInt(from.UnixNano()/int64(time.Millisecond), 10)},
			"to":   []string{strconv.FormatInt(to.UnixNano()/int64(time.Millisecond), 10)},
		}.Encode(),
	}
	resp, err := c.http.Get(u.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Get(ctx).Warnf("Grafana responded with status %d.", resp.StatusCode)
	}

	var response []annotation
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to decode JSON response")
	}
	for i, r := range response {
		r.decode()
		response[i] = r
	}
	return response, nil
}

type alertRule struct {
	Title        string             `json:"title"`
	Condition    string             `json:"conditions"`
	Data         []alertQuery       `json:"data"`
	UID          string             `json:"uid"`
	DashboardUID *string            `json:"dashboardUID"`
	PanelID      *int64             `json:"panelID"`
	RuleGroup    string             `json:"ruleGroup"`
	For          int64              `json:"for"`
	Annotations  map[string]*string `json:"annotations"`
	Labels       map[string]string  `json:"labels"`
	NoDataState  *string            `json:"noDataState"`
	ExecErrState *string            `json:"execErrState"`
}

type alertQuery struct {
	RefID             string `json:"refId"`
	QueryType         string `json:"queryType"`
	RelativeTimeRange struct {
		From time.Duration `json:"from"`
		To   time.Duration `json:"to"`
	} `json:"relativeTimeRange"`
	DatasourceUID string          `json:"datasourceUid"`
	Model         json.RawMessage `json:"model"`
}

type alertRuleParam struct {
	DatasourceUID string
	NamespaceUID  string
	Instance      string
}

func (c *Client) HealthAlertsEnabledMap(ctx context.Context, instances ...string) (map[string]bool, error) {
	enabledMap := make(map[string]bool)
	for _, instance := range instances {
		enabledMap[instance] = true
	}

	db, err := sql.Open("sqlite3", c.db)
	if err != nil {
		return nil, err
	}

	alertFiles, err := c.alertFiles()
	if err != nil {
		return nil, err
	}

	if len(alertFiles) == 0 || len(instances) == 0 {
		return map[string]bool{}, nil
	}

	for _, alertFile := range alertFiles {
		tpl, err := template.ParseFiles(alertFile)
		if err != nil {
			return nil, err
		}

		for _, instance := range instances {
			var alertRuleBytes bytes.Buffer
			if err = tpl.Execute(&alertRuleBytes, map[string]interface{}{
				"Instance": instance,
			}); err != nil {
				return nil, err
			}

			var rule alertRule
			if err = json.Unmarshal(alertRuleBytes.Bytes(), &rule); err != nil {
				return nil, err
			}
			if rule.UID == "" {
				continue
			}

			var id int64
			err = db.QueryRowContext(ctx, "SELECT id FROM alert_rule WHERE uid = ?", rule.UID).Scan(&id)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			if err == sql.ErrNoRows || id == 0 {
				enabledMap[instance] = false
			}
		}
	}

	return enabledMap, nil
}

func (c *Client) alertFiles() ([]string, error) {
	alertFiles := make([]string, 0)

	if err := filepath.WalkDir(c.alertsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if path == c.alertsPath {
			return nil
		}

		if d.IsDir() {
			return fs.SkipDir
		}

		if !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}

		alertFiles = append(alertFiles, path)
		return nil
	}); err != nil {
		return nil, err
	}

	return alertFiles, nil
}

func (c *Client) EnableHealthAlerts(ctx context.Context, instance string) error {
	alertFiles, err := c.alertFiles()
	if err != nil {
		return err
	}

	if len(alertFiles) == 0 {
		return nil
	}

	db, err := sql.Open("sqlite3", c.db)
	if err != nil {
		return err
	}

	param := alertRuleParam{
		Instance: instance,
	}
	err = db.QueryRowContext(ctx, "SELECT uid FROM data_source WHERE org_id = ? AND name = ?", defaultOrgID, defaultDatasource).Scan(&param.DatasourceUID)
	if err != nil {
		return err
	}
	err = db.QueryRowContext(ctx, "SELECT uid FROM dashboard WHERE org_id = ? AND folder_id = 0 AND title = ?", defaultOrgID, defaultAlertRuleNamespace).Scan(&param.NamespaceUID)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, alertFile := range alertFiles {
		tpl, err := template.ParseFiles(alertFile)
		if err != nil {
			return err
		}

		var alertRuleBytes bytes.Buffer
		if err = tpl.Execute(&alertRuleBytes, param); err != nil {
			return err
		}

		var rule alertRule
		if err = json.Unmarshal(alertRuleBytes.Bytes(), &rule); err != nil {
			return err
		}

		dataBytes, err := json.Marshal(rule.Data)
		if err != nil {
			return err
		}

		if rule.Annotations == nil {
			rule.Annotations = make(map[string]*string)
		}
		if _, exists := rule.Annotations["__dashboardUid__"]; !exists && rule.DashboardUID != nil {
			rule.Annotations["__dashboardUid__"] = rule.DashboardUID
		}
		if _, exists := rule.Annotations["__panelId__"]; !exists && rule.PanelID != nil {
			panelID := fmt.Sprintf("%d", *rule.PanelID)
			rule.Annotations["__panelId__"] = &panelID
		}
		annotationsBytes, err := json.Marshal(rule.Annotations)
		if err != nil {
			return err
		}

		labelsBytes, err := json.Marshal(rule.Labels)
		if err != nil {
			return err
		}

		ruleGroupIdx := 1
		if rule.RuleGroup != "" {
			var groupRuleCount int
			if err = db.QueryRowContext(ctx, "SELECT COUNT(1) FROM alert_rule WHERE org_id = ? AND namespace_uid = ? AND rule_group = ?", defaultOrgID, param.NamespaceUID, rule.RuleGroup).Scan(&groupRuleCount); err != nil {
				return err
			}
			ruleGroupIdx += groupRuleCount
		}

		alertRuleColumns := []string{"org_id", "title", "condition", "data", "updated", "version", "uid", "namespace_uid", "rule_group", "for", "annotations", "labels", "dashboard_uid", "panel_id", "rule_group_idx"}
		alertRuleArgs := []interface{}{
			defaultOrgID, rule.Title, rule.Condition, dataBytes,
			time.Now(), 1, rule.UID, param.NamespaceUID, rule.RuleGroup,
			rule.For, annotationsBytes, labelsBytes, rule.DashboardUID,
			rule.PanelID, ruleGroupIdx,
		}
		alertRuleVersionColumns := []string{"rule_org_id", "title", "condition", "data", "created", "version", "rule_uid", "rule_namespace_uid", "rule_group", "for", "annotations", "labels", "rule_group_idx", "parent_version", "restored_from", "interval_seconds"}
		alertRuleVersionArgs := []interface{}{
			defaultOrgID, rule.Title, rule.Condition, dataBytes,
			time.Now(), 1, rule.UID, param.NamespaceUID, rule.RuleGroup,
			rule.For, annotationsBytes, labelsBytes, ruleGroupIdx, 0, 0,
			defaultIntervalSeconds,
		}
		if rule.NoDataState != nil {
			alertRuleColumns = append(alertRuleColumns, "no_data_state")
			alertRuleVersionColumns = append(alertRuleVersionColumns, "no_data_state")
			alertRuleArgs = append(alertRuleArgs, *rule.NoDataState)
			alertRuleVersionArgs = append(alertRuleVersionArgs, *rule.NoDataState)
		}
		if rule.ExecErrState != nil {
			alertRuleColumns = append(alertRuleColumns, "exec_err_state")
			alertRuleVersionColumns = append(alertRuleVersionColumns, "exec_err_state")
			alertRuleArgs = append(alertRuleArgs, *rule.ExecErrState)
			alertRuleVersionArgs = append(alertRuleVersionArgs, *rule.ExecErrState)
		}

		if _, err = tx.ExecContext(
			ctx,
			fmt.Sprintf(`
				INSERT OR IGNORE INTO alert_rule (%s)
				VALUES (?%s)
			`, strings.Join(alertRuleColumns, ","), strings.Repeat(",?", len(alertRuleColumns)-1)),
			alertRuleArgs...,
		); err != nil {
			return err
		}

		if _, err = tx.ExecContext(
			ctx,
			fmt.Sprintf(`
				INSERT OR IGNORE INTO alert_rule_version (%s)
				VALUES (?%s)
			`, strings.Join(alertRuleVersionColumns, ","), strings.Repeat(",?", len(alertRuleVersionColumns)-1)),
			alertRuleVersionArgs...,
		); err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *Client) DisableHealthAlerts(ctx context.Context, instance string) error {
	alertFiles, err := c.alertFiles()
	if err != nil {
		return err
	}

	if len(alertFiles) == 0 {
		return nil
	}

	for _, alertFile := range alertFiles {
		tpl, err := template.ParseFiles(alertFile)
		if err != nil {
			return err
		}

		var alertRuleBytes bytes.Buffer
		if err = tpl.Execute(&alertRuleBytes, map[string]interface{}{
			"Instance": instance,
		}); err != nil {
			return err
		}

		var rule alertRule
		if err = json.Unmarshal(alertRuleBytes.Bytes(), &rule); err != nil {
			return err
		}
		if rule.UID == "" {
			continue
		}

		if err = c.deleteAlertRule(ctx, rule.UID); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) deleteAlertRule(ctx context.Context, uid string) error {
	u := &url.URL{
		Scheme: "http",
		Host:   c.addr,
		Path:   fmt.Sprintf("/api/v1/provisioning/alert-rules/%s", uid),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(os.Getenv("SERVER_USER"), os.Getenv("SERVER_PASSWORD"))
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.Errorf("got an unexpected status code from grafana's delete alert rule api: %d", resp.StatusCode)
	}

	return nil
}
