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
	"hash/fnv"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	prometheusapi "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shatteredsilicon/ssm-managed/utils"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
)

const (
	defaultOrgID              = 1
	defaultAlertRuleNamespace = "Alerts"
	defaultDatasource         = "Prometheus"
	defaultNoDataState        = "NoData"
	defaultExecErrState       = "Alerting"
)

const (
	AlertRuleStatusDisabled = iota
	AlertRuleStatusEnabled
	AlertRuleStatusPartiallyEnabled
)

const (
	AlertRuleCategorySystem  = "system"
	AlertRuleCategoryMySQL   = "mysql"
	AlertRuleCategoryMongoDB = "mongodb"
)

var tplFuncMap = template.FuncMap{
	"mul": func(args ...any) (float64, error) {
		val := float64(1)
		for _, arg := range args {
			switch v := arg.(type) {
			case int8:
			case int16:
			case int32:
			case int64:
			case float32:
			case float64:
				val *= float64(v)
			case nil:
				continue
			default:
				return 0, errors.Errorf("unsupported type %T for the mul function", v)
			}
		}
		return val, nil
	},
}

// Client represents a client for Grafana API.
type Client struct {
	addr       string
	db         string
	alertsPath string
	http       *http.Client
	prometheus prometheusapi.Client
}

// NewClient creates a new client for given Grafana address.
func NewClient(addr string, db string, alertsPath string, prometheus prometheusapi.Client) *Client {
	return &Client{
		addr:       addr,
		db:         db,
		alertsPath: alertsPath,
		http:       &http.Client{},
		prometheus: prometheus,
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
	OrgID        int64              `json:"orgID"`
	FolderUID    string             `json:"folderUID"`
	Title        string             `json:"title"`
	Condition    string             `json:"condition"`
	Data         []alertQuery       `json:"data"`
	UID          string             `json:"uid"`
	DashboardUID *string            `json:"dashboardUID"`
	PanelID      *int64             `json:"panelID"`
	RuleGroup    string             `json:"ruleGroup"`
	For          string             `json:"for"`
	Annotations  map[string]*string `json:"annotations"`
	Labels       map[string]string  `json:"labels"`
	NoDataState  string             `json:"noDataState"`
	ExecErrState string             `json:"execErrState"`
	Category     string             `json:"category"`
}

func (rule alertRule) checksumUID() string {
	uid := fnv.New64a()
	io.WriteString(uid, rule.UID)
	return fmt.Sprintf("%x", uid.Sum(nil))
}

type alertQuery struct {
	RefID             string `json:"refId"`
	QueryType         string `json:"queryType"`
	RelativeTimeRange struct {
		From time.Duration `json:"from"`
		To   time.Duration `json:"to"`
	} `json:"relativeTimeRange"`
	DatasourceUID string         `json:"datasourceUid"`
	Model         alertRuleModel `json:"model"`
}

type alertRuleModel struct {
	Expr       string `json:"expr"`
	Datasource struct {
		Type string `json:"type"`
		UID  string `json:"uid"`
	} `json:"datasource"`
	Hide          bool                 `json:"hide"`
	IntervalMs    int                  `json:"intervalMs"`
	MaxDataPoints int                  `json:"maxDataPoints"`
	RefID         string               `json:"refId"`
	LegendFormat  string               `json:"legendFormat"`
	Type          string               `json:"type"`
	Conditions    []alertRuleCondition `json:"conditions"`
}

type alertRuleCondition struct {
	Evaluator struct {
		Params []any  `json:"params"`
		Type   string `json:"type"`
	} `json:"evaluator"`
	Operator struct {
		Type string `json:"type"`
	} `json:"operator"`
	Query struct {
		Params []string `json:"params"`
	} `json:"query"`
	Reducer struct {
		Type   string   `json:"type"`
		Params []string `json:"params"`
	} `json:"reducer"`
	Type string `json:"type"`
}

func (c alertRuleCondition) MarshalJSON() ([]byte, error) {
	type condition alertRuleCondition

	for i := range c.Evaluator.Params {
		switch v := c.Evaluator.Params[i].(type) {
		case string:
			c.Evaluator.Params[i], _ = strconv.ParseFloat(v, 64)
		case int, int8, int16, int32, int64, float32, float64:
			continue
		default:
			return nil, errors.Errorf("unsupported param type %T for the condition evaluator", v)
		}
	}

	return json.Marshal(condition(c))
}

type alertRuleParam struct {
	DatasourceUID string
	NamespaceUID  string
	Instance      string
	InitialValue  float64
}

func (c *Client) HealthAlertsStateMap(ctx context.Context, instances ...string) (map[string]map[string]int32, error) {
	stateMap := make(map[string]map[string]int32)

	db, err := sql.Open("sqlite3", c.db)
	if err != nil {
		return nil, err
	}

	alertFiles, err := c.alertFiles()
	if err != nil {
		return nil, err
	}

	if len(alertFiles) == 0 || len(instances) == 0 {
		return map[string]map[string]int32{}, nil
	}

	for _, instance := range instances {
		for _, alertFile := range alertFiles {
			tpl, err := template.New(path.Base(alertFile)).Funcs(tplFuncMap).ParseFiles(alertFile)
			if err != nil {
				return nil, err
			}

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
			err = db.QueryRowContext(ctx, "SELECT id FROM alert_rule WHERE uid = ? OR uid = ?", rule.UID, rule.checksumUID()).Scan(&id)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			state, stateExists := stateMap[instance]
			if !stateExists {
				state = make(map[string]int32)
			}

			if err == sql.ErrNoRows || id == 0 {
				if state[rule.Category] == AlertRuleStatusEnabled {
					state[rule.Category] = AlertRuleStatusPartiallyEnabled
				} else if _, ok := state[rule.Category]; !ok {
					state[rule.Category] = AlertRuleStatusDisabled
				}
			} else {
				if _, ok := state[rule.Category]; !ok {
					state[rule.Category] = AlertRuleStatusEnabled
				} else if state[rule.Category] == AlertRuleStatusDisabled {
					state[rule.Category] = AlertRuleStatusPartiallyEnabled
				}
			}

			stateMap[instance] = state
		}
	}

	return stateMap, nil
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

func (c *Client) EnableHealthAlerts(ctx context.Context, instance string, categories ...string) error {
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

	defer func() {
		if err != nil {
			go c.DisableHealthAlerts(context.TODO(), instance, categories...)
		}
	}()

	for _, alertFile := range alertFiles {
		var tpl *template.Template
		tpl, err = template.New(path.Base(alertFile)).Funcs(tplFuncMap).ParseFiles(alertFile)
		if err != nil {
			return err
		}

		if rawRuleString := tpl.Root.String(); strings.Contains(rawRuleString, ".InitialValue") {
			var tmpRule alertRule
			if err = json.Unmarshal([]byte(rawRuleString), &tmpRule); err != nil {
				return err
			}

			if len(categories) > 0 && !utils.SliceContains(categories, tmpRule.Category) {
				continue
			}

			for _, d := range tmpRule.Data {
				if d.Model.Datasource.Type != "prometheus" || d.Model.Expr == "" {
					continue
				}

				exprTpl, err := template.New("").Funcs(tplFuncMap).Parse(d.Model.Expr)
				if err != nil {
					return err
				}

				var queryBytes bytes.Buffer
				if err = exprTpl.Execute(&queryBytes, param); err != nil {
					return err
				}

				initialVal, _, err := v1.NewAPI(c.prometheus).Query(ctx, queryBytes.String(), time.Now())
				if err != nil {
					return err
				}

				switch initialVal.Type() {
				case model.ValScalar:
					param.InitialValue = float64(initialVal.(*model.Scalar).Value)
				case model.ValMatrix:
					matrix := initialVal.(model.Matrix)
					if len(matrix) > 0 && len(matrix[0].Values) > 0 {
						param.InitialValue = float64(matrix[0].Values[0].Value)
					}
				case model.ValVector:
					vector := initialVal.(model.Vector)
					if len(vector) > 0 {
						param.InitialValue = float64(vector[0].Value)
					}
				}
			}
		}

		var alertRuleBytes bytes.Buffer
		if err = tpl.Execute(&alertRuleBytes, param); err != nil {
			return err
		}

		var rule alertRule
		if err = json.Unmarshal(alertRuleBytes.Bytes(), &rule); err != nil {
			return err
		}

		if rule.UID == "" {
			continue
		}

		if len(categories) > 0 && !utils.SliceContains(categories, rule.Category) {
			continue
		}

		// make sure uid is less than 40 characters
		rule.UID = rule.checksumUID()

		if rule.OrgID == 0 {
			rule.OrgID = defaultOrgID
		}

		if rule.FolderUID == "" {
			rule.FolderUID = param.NamespaceUID
		}

		if rule.RuleGroup == "" {
			rule.RuleGroup = param.Instance
		}

		if rule.NoDataState == "" {
			rule.NoDataState = defaultNoDataState
		}
		if rule.ExecErrState == "" {
			rule.ExecErrState = defaultExecErrState
		}

		if rule.Annotations == nil {
			rule.Annotations = make(map[string]*string)
		}
		if rule.Annotations["__dashboardUid__"] == nil && rule.Annotations["__panelId__"] == nil && rule.DashboardUID != nil && rule.PanelID != nil {
			rule.Annotations["__dashboardUid__"] = rule.DashboardUID
			dashboardURL := fmt.Sprintf("http/../d/%s?var-host=%s", *rule.DashboardUID, instance)
			rule.Annotations["Dashboard URL"] = &dashboardURL

			panelID := fmt.Sprintf("%d", *rule.PanelID)
			rule.Annotations["__panelId__"] = &panelID
			panelURL := fmt.Sprintf("http/../d/%s?var-host=%s&viewPanel=panel-%d", *rule.DashboardUID, instance, *rule.PanelID)
			rule.Annotations["Panel URL"] = &panelURL
		}

		body, _ := json.Marshal(rule)
		if err = c.postAlertRule(ctx, body); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) DisableHealthAlerts(ctx context.Context, instance string, categories ...string) error {
	alertFiles, err := c.alertFiles()
	if err != nil {
		return err
	}

	if len(alertFiles) == 0 {
		return nil
	}

	for _, alertFile := range alertFiles {
		tpl, err := template.New(path.Base(alertFile)).Funcs(tplFuncMap).ParseFiles(alertFile)
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

		if len(categories) > 0 && !utils.SliceContains(categories, rule.Category) {
			continue
		}

		if err = c.deleteAlertRule(ctx, rule.UID); err != nil {
			return err
		}

		if err = c.deleteAlertRule(ctx, rule.checksumUID()); err != nil {
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

func (c *Client) postAlertRule(ctx context.Context, body []byte) error {
	u := &url.URL{
		Scheme: "http",
		Host:   c.addr,
		Path:   "/api/v1/provisioning/alert-rules",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Disable-Provenance", "true") // enable editing

	req.SetBasicAuth(os.Getenv("SERVER_USER"), os.Getenv("SERVER_PASSWORD"))
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode < 200 || resp.StatusCode > 299) && resp.StatusCode != http.StatusConflict {
		return errors.Errorf("got an unexpected status code from grafana's post alert rule api: %d", resp.StatusCode)
	}

	return nil
}
