package prometheus

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
	"github.com/shatteredsilicon/ssm-managed/models"
)

const (
	targetsURI = "api/v1/targets"
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
		if target.Labels.Instance == "" || target.Labels.Instance == string(models.PMMServerNodeType) {
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
