package prometheus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/golang/snappy"
	"github.com/pkg/errors"
	"github.com/shatteredsilicon/ssm-managed/api"
	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	targetsURI     = "api/v1/targets"
	labelValuesURI = "api/v1/label/%s/values"
	remoteWriteURI = "api/v1/write"
)

var scrapePoolServiceMap = map[string]models.AgentType{
	"linux":             models.ClientNodeExporterAgentType,
	"mongodb":           models.ClientMongoDBExporterAgentType,
	"proxysql":          models.ClientProxySQLExporterAgentType,
	"remote-proxysql":   models.ProxySQLExporterAgentType,
	"postgresql":        models.ClientPostgresExporterAgentType,
	"remote-postgresql": models.PostgresExporterAgentType,
	"mysql-hr":          models.ClientMySQLdExporterAgentType,
	"mysql-mr":          models.ClientMySQLdExporterAgentType,
	"mysql-lr":          models.ClientMySQLdExporterAgentType,
	"remote-mysql-hr":   models.MySQLdExporterAgentType,
	"remote-mysql-mr":   models.MySQLdExporterAgentType,
	"remote-mysql-lr":   models.MySQLdExporterAgentType,
	"rds-mysql-hr":      models.MySQLdExporterAgentType,
	"rds-mysql-mr":      models.MySQLdExporterAgentType,
	"rds-mysql-lr":      models.MySQLdExporterAgentType,
	"rds-basic":         models.RDSExporterAgentType,
	"rds-enhanced":      models.RDSExporterAgentType,
	"snmp":              models.SNMPExporterAgentType,
}

// TargetActiveTarget active target structure of prometheus GET targets api
type TargetActiveTarget struct {
	DiscoveredLabels struct {
		Address string `json:"__address__"`
	} `json:"discoveredLabels"`
	Labels struct {
		Instance string `json:"instance"`
		Job      string `json:"job"`
	} `json:"labels"`
	ScrapePool string `json:"scrapePool"`
}

// TargetData data structure of prometheus GET targets api
type TargetData struct {
	ActiveTargets []TargetActiveTarget `json:"activeTargets"`
}

// TargetResponse response structrue of prometheus GET targets api
type TargetResponse struct {
	Status string     `json:"status"`
	Data   TargetData `json:"data"`
}

// NodeService service of node
type NodeService struct {
	Name     string
	Type     models.AgentType
	Endpoint string
}

// GetNodeServices returns services of ndoe
func (svc *Service) GetNodeServices(ctx context.Context) ([]NodeService, error) {
	u := *svc.baseURL
	u.Path = path.Join(u.Path, targetsURI)
	resp, err := svc.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.Errorf("unexpected api %s returns: code: %d, data: %s", u.String(), resp.StatusCode, string(b))
	}

	var data TargetResponse
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	if data.Status != "success" {
		return nil, errors.Errorf("unexpected api %s status: %s", u.String(), data.Status)
	}

	services := make([]NodeService, 0)
	for _, target := range data.Data.ActiveTargets {
		if target.Labels.Instance == "" || target.Labels.Instance == string(models.SSMServerNodeType) {
			continue
		}

		agentType, ok := scrapePoolServiceMap[target.ScrapePool]
		if !ok {
			continue
		}

		exists := false
		for _, service := range services {
			if service.Name == target.Labels.Instance && service.Type == agentType {
				exists = true
				break
			}
		}

		if !exists {
			services = append(services, NodeService{
				Name:     target.Labels.Instance,
				Type:     agentType,
				Endpoint: target.DiscoveredLabels.Address,
			})
		}
	}

	return services, nil
}

// LabelValues response struct of prometheus' label values API
type LabelValues struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

func (svc *Service) GetLabelValues(label string) (*LabelValues, error) {
	u := *svc.baseURL
	u.Path = path.Join(u.Path, fmt.Sprintf(labelValuesURI, label))
	resp, err := svc.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.Errorf("unexpected api %s returns: code: %d, data: %s", u.String(), resp.StatusCode, string(b))
	}

	var data LabelValues
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	if data.Status != "success" {
		return nil, errors.Errorf("unexpected api %s status: %s", u.String(), data.Status)
	}

	return &data, nil
}

func (svc *Service) RemoveNode(ctx context.Context, nodeID string) error {
	deleteSeries := func(ctx context.Context, nodeName string) error {
		activeTargets, err := svc.GetNodeServices(ctx)
		if err != nil {
			return fmt.Errorf("get active targets from prometheus for %s failed: %s", nodeName, err.Error())
		}

		for _, target := range activeTargets {
			if string(target.Name) == nodeName { // service has been re-added
				return nil
			}
		}

		err = svc.DeleteSeries(map[string]string{
			"instance=": nodeName,
		})
		if err != nil {
			return fmt.Errorf("delete metrics data for %s failed: %s", nodeName, err.Error())
		}

		return nil
	}

	deleteSeries(ctx, nodeID)

	// continually remove prometheus data incase there are some ongoing metrics task
	retryTimes := 35
	for i := 0; i < retryTimes; i++ {
		if retryTimes < 30 {
			<-time.NewTimer(1 * time.Second).C // for second-level jobs
		} else {
			<-time.NewTimer(1 * time.Minute).C // for minute-level jobs
		}

		err := deleteSeries(ctx, nodeID)
		if err != nil {
			logrus.Errorf("delete metrics data for %s failed: %s, try %d", nodeID, err.Error(), i+1)
			continue
		}

		err = svc.CleanTombstones(context.Background())
		if err != nil {
			logrus.Errorf("clean tombstones for %s failed: %s, try %d", nodeID, err.Error(), i+1)
			continue
		}
	}

	return fmt.Errorf("delete metrics data for %s failed, tried %d times", nodeID, retryTimes)
}

// WriteMetrics calls remote write API of prometheus
func (svc *Service) WriteMetrics(ctx context.Context, payload api.WriteRequest) error {
	buf, err := proto.Marshal(&payload)
	if err != nil {
		return err
	}

	u := *svc.baseURL
	u.Path = path.Join(u.Path, remoteWriteURI)
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), bytes.NewBuffer(snappy.Encode(nil, buf)))
	if err != nil {
		return err
	}

	// following headers are required as described
	// in prometheus remote write spec ->
	// https://prometheus.io/docs/concepts/remote_write_spec/
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	postRequest := func() (int, []byte, error) {
		resp, err := svc.client.Do(req)
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		return resp.StatusCode, body, err
	}

	retryTimes := 10
	for i := 0; i < retryTimes; i++ {
		statusCode, body, err := postRequest()
		if err != nil {
			return err
		}

		if statusCode >= 200 && statusCode < 300 {
			return nil
		}

		if statusCode >= 400 && statusCode < 500 && statusCode != http.StatusTooManyRequests {
			return errors.New(fmt.Sprintf("(%d): %s", statusCode, string(body)))
		}

		logger.Get(ctx).WithField("component", "prometheus").Warnf("remote write api returned an (%d): %s response, tried %d times", statusCode, body, i+1)
		time.Sleep(time.Second)
	}

	return errors.New(fmt.Sprintf("remote write api didn't response with an success status code, tried %d times", retryTimes))
}
