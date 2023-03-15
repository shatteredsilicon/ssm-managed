package snmp

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"sort"
	"strconv"

	"github.com/go-sql-driver/mysql"
	servicelib "github.com/percona/kardianos-service"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"gopkg.in/yaml.v2"

	"github.com/AlekSi/pointer"
	prometheusModel "github.com/prometheus/common/model"
	snmpConfig "github.com/shatteredsilicon/snmp_exporter/config"
	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/services"
	"github.com/shatteredsilicon/ssm-managed/services/prometheus"
	"github.com/shatteredsilicon/ssm-managed/utils"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
	"github.com/shatteredsilicon/ssm-managed/utils/ports"
)

const (
	defaultSNMPPort uint32 = 161
)

var (
	snmpEngine   = "SNMP"
	snmpModule   = "ssm_mib"
	snmpVersions = []string{"1", "2", "3"}
)

type ServiceConfig struct {
	SNMPExporterPath  string
	SNMPGeneratorPath string
	SNMPConfigDir     string

	Prometheus    *prometheus.Service
	Supervisor    services.Supervisor
	DB            *reform.DB
	PortsRegistry *ports.Registry
}

// Service is responsible for interactions with SNMP.
type Service struct {
	*ServiceConfig
	httpClient    *http.Client
	ssmServerNode *models.Node
}

func (svc *Service) generatorConfigPath(name string) string {
	return path.Join(svc.SNMPConfigDir, name+"-generator.yml")
}

func (svc *Service) baseGeneratorConfigPath() string {
	return path.Join(svc.SNMPConfigDir, "generator.yml")
}

func (svc *Service) exporterConfigPath(name string) string {
	return path.Join(svc.SNMPConfigDir, name+"-snmp.yml")
}

// NewService creates a new service.
func NewService(config *ServiceConfig) (*Service, error) {
	var node models.Node
	err := config.DB.FindOneTo(&node, "type", models.PMMServerNodeType)
	if err != nil {
		return nil, err
	}

	for _, path := range []*string{
		&config.SNMPExporterPath,
	} {
		if *path == "" {
			continue
		}
		p, err := exec.LookPath(*path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		*path = p
	}

	svc := &Service{
		ServiceConfig: config,
		httpClient:    new(http.Client),
		ssmServerNode: &node,
	}
	return svc, nil
}

func (svc *Service) ApplyPrometheusConfiguration(ctx context.Context, q *reform.Querier) error {
	snmp := &prometheus.ScrapeConfig{
		JobName:        "snmp",
		ScrapeInterval: "1s",
		ScrapeTimeout:  "1s",
		MetricsPath:    "/snmp",
		Params: url.Values{
			"module": []string{snmpModule},
		},
		HonorLabels: true,
		RelabelConfigs: []prometheus.RelabelConfig{
			{
				SourceLabels: prometheusModel.LabelNames{"__address__"},
				TargetLabel:  "__param_target",
			},
			{
				SourceLabels: prometheusModel.LabelNames{"__param_target"},
				TargetLabel:  "host",
			},
		},
	}

	nodes, err := q.FindAllFrom(models.RemoteNodeTable, "type", models.RemoteNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.RemoteNode)

		var service models.SNMPService
		if e := q.SelectOneTo(&service, "WHERE node_id = ? AND type = ?", node.ID, models.SNMPServiceType); e == sql.ErrNoRows {
			continue
		} else if e != nil {
			return errors.WithStack(e)
		}

		agents, err := models.AgentsForServiceID(q, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.SNMPExporterAgentType:
				a := models.SNMPExporter{ID: agent.ID}
				if e := q.Reload(&a); e != nil {
					return errors.WithStack(e)
				}
				logger.Get(ctx).WithField("component", "snmp").Infof("%s %s %s %d", a.Type, node.Name, node.Region, *a.ListenPort)

				addressExists := false
				for _, rc := range snmp.RelabelConfigs {
					if rc.TargetLabel == "__address__" {
						addressExists = true
						break
					}
				}
				if !addressExists {
					snmp.RelabelConfigs = append(snmp.RelabelConfigs, prometheus.RelabelConfig{
						TargetLabel: "__address__",
						Replacement: fmt.Sprintf("127.0.0.1:%d", *a.ListenPort),
					})
				}

				sc := prometheus.StaticConfig{
					Targets: []string{fmt.Sprintf("%s:%d", *service.Address, *service.Port)},
					Labels: []prometheus.LabelPair{
						{Name: "instance", Value: node.Name},
						{Name: "real_job", Value: "snmp"},
					},
				}
				snmp.StaticConfigs = append(snmp.StaticConfigs, sc)
			}
		}
	}

	// sort by region and name
	sorterFor := func(sc []prometheus.StaticConfig) func(int, int) bool {
		return func(i, j int) bool {
			if sc[i].Labels[0].Value != sc[j].Labels[0].Value {
				return sc[i].Labels[0].Value < sc[j].Labels[0].Value
			}
			return sc[i].Labels[1].Value < sc[j].Labels[1].Value
		}
	}
	sort.Slice(snmp.StaticConfigs, sorterFor(snmp.StaticConfigs))

	return svc.Prometheus.SetScrapeConfigs(ctx, false, snmp)
}

func (svc *Service) snmpExporterServiceConfig(agent *models.SNMPExporter, instanceName string) *servicelib.Config {
	name := models.NameForSupervisor(agent.Type, *agent.ListenPort)

	return &servicelib.Config{
		Name:        name,
		DisplayName: name,
		Description: name,
		Executable:  svc.SNMPExporterPath,
		Arguments: []string{
			fmt.Sprintf("--config.file=%s", svc.exporterConfigPath(instanceName)),
			fmt.Sprintf("--web.listen-address=127.0.0.1:%d", *agent.ListenPort),
		},
	}
}

func (svc *Service) addSNMPExporter(ctx context.Context, tx *reform.TX, service *models.SNMPService, instanceName string) error {
	port, err := svc.PortsRegistry.Reserve()
	if err != nil {
		return err
	}

	// insert snmp_exporter agent and association
	agent := &models.SNMPExporter{
		Type:         models.SNMPExporterAgentType,
		RunsOnNodeID: svc.ssmServerNode.ID,

		ListenPort: &port,
	}
	if err := tx.Insert(agent); err != nil {
		return errors.WithStack(err)
	}
	if err := tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID}); err != nil {
		return errors.WithStack(err)
	}

	// start snmp_exporter agent
	if svc.SNMPExporterPath != "" {
		cfg := svc.snmpExporterServiceConfig(agent, instanceName)
		if err := svc.Supervisor.Start(ctx, cfg); err != nil {
			return err
		}
	}

	return nil
}

// Add adds SNMP instance
func (svc *Service) Add(
	ctx context.Context,
	name, address string,
	port uint32,
	username, password string,
	version, community string,
	securityLevel string,
	authProtocol, privProtocol string,
	privPassword string,
	contextName string,
) (int32, error) {
	if name == "" {
		return 0, status.Error(codes.InvalidArgument, "SNMP instance name is not given.")
	}
	if !utils.SliceContains(snmpVersions, version) {
		return 0, status.Error(codes.InvalidArgument, "Unsupported SNMP version")
	}

	if port == 0 {
		port = defaultSNMPPort
	}

	var id int32
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		// insert node
		node := &models.RemoteNode{
			Type:   models.RemoteNodeType,
			Name:   name,
			Region: string(models.RemoteNodeRegion),
		}
		if err := tx.Insert(node); err != nil {
			if err, ok := err.(*mysql.MySQLError); ok && err.Number == 0x426 {
				return status.Errorf(codes.AlreadyExists, "SNMP instance %q already exists.", node.Name)
			}
			return errors.WithStack(err)
		}
		id = node.ID

		// insert service
		service := &models.SNMPService{
			Type:   models.SNMPServiceType,
			NodeID: node.ID,

			Address:       &address,
			Port:          pointer.ToUint16(uint16(port)),
			Engine:        &snmpEngine,
			EngineVersion: &version,
		}
		if err := tx.Insert(service); err != nil {
			return errors.WithStack(err)
		}

		intVersion, _ := strconv.Atoi(version)
		if err := svc.generateSNMPConfig(
			name, username, password, intVersion, community, securityLevel,
			authProtocol, privProtocol, privPassword, contextName,
		); err != nil {
			return err
		}

		if err := svc.addSNMPExporter(ctx, tx, service, name); err != nil {
			return err
		}

		return svc.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})

	return id, err
}

func (svc *Service) generateSNMPConfig(
	name string,
	username, password string,
	version int, community string,
	securityLevel string,
	authProtocol, privProtocol string,
	privPassword string,
	contextName string,
) error {
	configBytes, err := ioutil.ReadFile(svc.baseGeneratorConfigPath())
	if err != nil {
		return err
	}

	var config snmpGeneratorConfig
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return err
	}

	if _, ok := config.Modules[snmpModule]; !ok {
		return errors.Errorf("snmp_exporter generator config file is broken")
	}

	config.Modules[snmpModule].WalkParams.Version = version
	config.Modules[snmpModule].WalkParams.Auth.Community = snmpConfig.Secret(community)
	config.Modules[snmpModule].WalkParams.Auth.AuthProtocol = authProtocol
	config.Modules[snmpModule].WalkParams.Auth.Username = username
	config.Modules[snmpModule].WalkParams.Auth.Password = snmpConfig.Secret(password)
	config.Modules[snmpModule].WalkParams.Auth.SecurityLevel = securityLevel
	config.Modules[snmpModule].WalkParams.Auth.PrivProtocol = privProtocol
	config.Modules[snmpModule].WalkParams.Auth.PrivPassword = snmpConfig.Secret(privPassword)
	config.Modules[snmpModule].WalkParams.Auth.ContextName = contextName

	// update snmp_exporter configuration
	snmpConfig.DoNotHideSecrets = true
	b, err := yaml.Marshal(config)
	snmpConfig.DoNotHideSecrets = false
	if err != nil {
		return err
	}
	b = append([]byte("# Managed by ssm-managed. DO NOT EDIT.\n"), b...)
	if err = ioutil.WriteFile(svc.generatorConfigPath(name), b, 0666); err != nil {
		return err
	}

	cmd := exec.Command(svc.SNMPGeneratorPath, "generate", "-i", svc.generatorConfigPath(name), "-o", svc.exporterConfigPath(name))
	return cmd.Run()
}
