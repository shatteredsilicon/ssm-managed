package snmp

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	"github.com/shatteredsilicon/ssm-managed/services/consul"
	"github.com/shatteredsilicon/ssm-managed/services/prometheus"
	"github.com/shatteredsilicon/ssm-managed/utils"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
	"github.com/shatteredsilicon/ssm-managed/utils/ports"
)

const (
	defaultSNMPPort      uint32 = 161
	defaultCommunity     string = "public"
	defaultVersion       string = "2"
	defaultUsername      string = "admin"
	defaultSecurityLevel string = "noAuthNoPriv"
)

var (
	snmpEngine   = "SNMP"
	snmpModule   = "ssm_mib"
	snmpAuth     = "ssm_auth"
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
	Consul        *consul.Client
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
	err := config.DB.FindOneTo(&node, "type", models.SSMServerNodeType)
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
			"auth":   []string{snmpAuth},
		},
		HonorLabels: true,
		RelabelConfigs: []prometheus.RelabelConfig{
			{
				SourceLabels: prometheusModel.LabelNames{"__address__"},
				TargetLabel:  "__param_target",
			},
			{
				TargetLabel: "job",
				Replacement: "linux",
			},
		},
		MetricRelabelConfigs: []prometheus.RelabelConfig{
			{
				SourceLabels: prometheusModel.LabelNames{"__name__"},
				TargetLabel:  "__name__",
				Regex:        "(node_disk_reads_completed|node_disk_writes_completed|node_disk_io_time_ms|node_intr|node_forks|node_context_switches|node_vmstat_pswpin|node_vmstat_pswpout|node_vmstat_pgpgin|node_vmstat_pgpgout|node_disk_write_time_ms|node_disk_read_time_ms|node_network_receive_bytes|node_network_transmit_bytes|node_disk_reads_merged|node_disk_writes_merged|node_disk_sectors_read|node_disk_sectors_written|node_netstat_Tcp_RetransSegs|node_netstat_Tcp_OutSegs|node_netstat_Tcp_InSegs|node_netstat_Udp_InDatagrams|node_netstat_Udp_InErrors|node_netstat_Udp_OutDatagrams|node_netstat_Udp_NoPorts|node_netstat_Udp_InCsumErrors|node_netstat_Udp_RcvbufErrors|node_netstat_Udp_SndbufErrors|node_netstat_UdpLite_InDatagrams|node_netstat_UdpLite_OutDatagrams|node_netstat_UdpLite_InCsumErrors|node_netstat_UdpLite_InErrors|node_netstat_UdpLite_RcvbufErrors|node_netstat_UdpLite_SndbufErrors|node_netstat_UdpLite_NoPorts|node_netstat_Icmp_InErrors|node_netstat_Icmp_OutErrors|node_netstat_Icmp_InDestUnreachs|node_netstat_Icmp_OutDestUnreachs|node_netstat_IcmpMsg_InType3|node_netstat_IcmpMsg_OutType3|node_netstat_Icmp_InCsumErrors|node_netstat_Icmp_InTimeExcds|node_netstat_Icmp_InMsgs|node_netstat_Icmp_InRedirects|node_netstat_Icmp_OutMsgs|node_netstat_Icmp_OutRedirects|node_netstat_Icmp_InEchoReps|node_netstat_Icmp_InEchos|node_netstat_Icmp_OutEchoReps|node_netstat_Icmp_OutEchos|node_netstat_Icmp_InAddrMaskReps|node_netstat_Icmp_OutAddrMasks|node_netstat_Icmp_InTimestamps|node_netstat_Icmp_OutTimestamps|node_netstat_Icmp_OutTimestampReps|node_netstat_Icmp_InTimestampReps|node_netstat_Icmp_InAddrMasks|node_netstat_Icmp_OutAddrMaskReps|node_netstat_Tcp_OutRsts|node_network_receive_packets|node_network_transmit_packets|node_network_receive_errs|node_network_transmit_errs|node_network_receive_drop|node_network_transmit_drop|node_network_receive_multicast|node_network_transmit_multicast|node_netstat_Tcp_InCsumErrors|node_netstat_Tcp_InErrs)",
				Replacement:  "${1}_total",
			},
			{
				SourceLabels: prometheusModel.LabelNames{"__name__"},
				TargetLabel:  "__name__",
				Regex:        "node_disk_bytes_read",
				Replacement:  "node_disk_read_bytes_total",
			},
			{
				SourceLabels: prometheusModel.LabelNames{"__name__"},
				TargetLabel:  "__name__",
				Regex:        "node_disk_bytes_written",
				Replacement:  "node_disk_written_bytes_total",
			},
			{
				SourceLabels: prometheusModel.LabelNames{"__name__"},
				TargetLabel:  "__name__",
				Regex:        "node_cpu",
				Replacement:  "node_cpu_seconds_total",
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

func (svc *Service) addSNMPExporter(ctx context.Context, tx *reform.TX, service *models.SNMPService, instanceName, username, password string) error {
	port, err := svc.PortsRegistry.Reserve()
	if err != nil {
		return err
	}

	// insert snmp_exporter agent and association
	agent := &models.SNMPExporter{
		Type:            models.SNMPExporterAgentType,
		RunsOnNodeID:    svc.ssmServerNode.ID,
		ListenPort:      &port,
		ServiceUsername: &username,
		ServicePassword: &password,
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
	if community == "" {
		community = defaultCommunity
	}

	// check if instance is added on client side
	added, err := svc.clientInstanceAdded(ctx, name)
	if err != nil {
		return 0, err
	}
	if added {
		return 0, status.Error(codes.AlreadyExists, fmt.Sprintf("Linux instance with name %s is already added", name))
	}

	if port == 0 {
		port = defaultSNMPPort
	}

	var id int32
	err = svc.DB.InTransaction(func(tx *reform.TX) error {
		// insert node
		node := &models.RemoteNode{
			Type:   models.RemoteNodeType,
			Name:   name,
			Region: string(models.RemoteNodeRegion),
		}
		if err := tx.Insert(node); err != nil {
			if err, ok := err.(*mysql.MySQLError); !ok || err.Number != 0x426 {
				return errors.WithStack(err)
			}

			err = tx.SelectOneTo(node, "WHERE type = ? AND name = ? AND region = ?", models.RemoteNodeType, name, string(models.RemoteNodeRegion))
			if err != nil {
				return errors.WithStack(err)
			}
		}
		id = node.ID

		// insert service
		service := &models.SNMPService{
			Type:   models.SNMPServiceType,
			NodeID: node.ID,

			PrivPassword: &privPassword,
			Address:      &address,
			Port:         pointer.ToUint16(uint16(port)),
			Engine:       &snmpEngine,
			EngineVersion: &models.SNMPEngineVersion{
				Version:       version,
				Community:     community,
				SecurityLevel: securityLevel,
				AuthProtocol:  authProtocol,
				PrivProtocol:  privPassword,
				ContextName:   contextName,
			},
		}
		if err := tx.Insert(service); err != nil {
			return errors.WithStack(err)
		}

		if err := svc.generateSNMPConfig(
			name, username, password, version, community, securityLevel,
			authProtocol, privProtocol, privPassword, contextName,
		); err != nil {
			return err
		}

		if err := svc.addSNMPExporter(ctx, tx, service, name, username, password); err != nil {
			return err
		}

		return svc.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})

	return id, err
}

func (svc *Service) generateSNMPConfig(
	name string,
	username, password string,
	version, community string,
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

	if config.Auths == nil {
		config.Auths = make(map[string]*snmpConfig.Auth)
	}
	intVersion, _ := strconv.Atoi(version)
	config.Auths[snmpAuth] = &snmpConfig.Auth{
		Community:     snmpConfig.Secret(community),
		SecurityLevel: securityLevel,
		Username:      username,
		Password:      snmpConfig.Secret(password),
		AuthProtocol:  authProtocol,
		PrivProtocol:  privProtocol,
		PrivPassword:  snmpConfig.Secret(privPassword),
		ContextName:   contextName,
		Version:       intVersion,
	}

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

	cmd := exec.Command(svc.SNMPGeneratorPath, "generate", "-g", svc.generatorConfigPath(name), "-o", svc.exporterConfigPath(name))
	return cmd.Run()
}

// Instance snmp instance
type Instance struct {
	Node    models.RemoteNode
	Service models.SNMPService
}

// List returns all snmp instances
func (svc *Service) List(ctx context.Context) ([]Instance, error) {
	var res []Instance
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		structs, e := tx.SelectAllFrom(models.RemoteNodeTable, "WHERE type = ? ORDER BY id", models.RemoteNodeType)
		if e != nil {
			return e
		}
		nodes := make([]models.RemoteNode, len(structs))
		for i, str := range structs {
			nodes[i] = *str.(*models.RemoteNode)
		}

		structs, e = tx.SelectAllFrom(models.SNMPServiceTable, "WHERE type = ? ORDER BY id", models.SNMPServiceType)
		if e != nil {
			return e
		}
		services := make([]models.SNMPService, len(structs))
		for i, str := range structs {
			services[i] = *str.(*models.SNMPService)
		}

		for _, node := range nodes {
			for _, service := range services {
				if node.ID == service.NodeID {
					res = append(res, Instance{
						Node:    node,
						Service: service,
					})
				}
			}
		}
		return nil
	})
	return res, err
}

func (svc *Service) clientInstanceAdded(ctx context.Context, name string) (bool, error) {
	node, err := svc.Consul.GetNode(name)
	if err != nil {
		logger.Get(ctx).Errorf("get consul services from node failed: %+v", err)
		return false, err
	}

	if node == nil {
		return false, nil
	}

	for _, service := range node.Services {
		t := models.AgentType(service.Service)
		if t == models.ClientNodeExporterAgentType {
			// instance is added on client side
			return true, nil
		}
	}

	return false, nil
}

// Restore configuration from database.
func (svc *Service) Restore(ctx context.Context, tx *reform.TX) error {
	nodes, err := tx.FindAllFrom(models.RemoteNodeTable, "type", models.RemoteNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.RemoteNode)

		snmpServices, e := tx.SelectAllFrom(models.SNMPServiceTable, "WHERE node_id = ? AND type = ?", node.ID, models.SNMPServiceType)
		if e != nil {
			return errors.WithStack(e)
		}
		if len(snmpServices) == 0 {
			continue
		}

		service := snmpServices[0].(*models.SNMPService)
		if service.EngineVersion == nil {
			service.EngineVersion = &models.SNMPEngineVersion{}
		}
		if service.EngineVersion.Version == "" {
			service.EngineVersion.Version = defaultVersion
		}
		if service.EngineVersion.Community == "" {
			service.EngineVersion.Community = defaultCommunity
		}
		if service.EngineVersion.SecurityLevel == "" {
			service.EngineVersion.SecurityLevel = defaultSecurityLevel
		}

		agents, err := models.AgentsForServiceID(tx.Querier, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.SNMPExporterAgentType:
				a := &models.SNMPExporter{ID: agent.ID}
				if err = tx.Reload(a); err != nil {
					return errors.WithStack(err)
				}

				username := defaultUsername
				var authPassword, privPasssword string
				if a.ServiceUsername != nil {
					username = *a.ServiceUsername
				}
				if a.ServicePassword != nil {
					authPassword = *a.ServicePassword
				}
				if service.PrivPassword != nil {
					privPasssword = *service.PrivPassword
				}
				if svc.SNMPGeneratorPath != "" {
					// generate snmp.yml if not exists
					if _, err := os.Stat(svc.exporterConfigPath(node.Name)); errors.Is(err, os.ErrNotExist) {
						if err := svc.generateSNMPConfig(
							node.Name, username, authPassword, service.EngineVersion.Version, service.EngineVersion.Community,
							service.EngineVersion.SecurityLevel, service.EngineVersion.AuthProtocol, service.EngineVersion.PrivProtocol,
							privPasssword, service.EngineVersion.ContextName,
						); err != nil {
							return err
						}
					}
				}
				if svc.SNMPExporterPath != "" {
					name := models.NameForSupervisor(a.Type, *a.ListenPort)

					err := svc.Supervisor.Status(ctx, name)
					if err == nil {
						if err = svc.Supervisor.Stop(ctx, name); err != nil {
							return err
						}
					}

					cfg := svc.snmpExporterServiceConfig(a, node.Name)
					if err = svc.Supervisor.Start(ctx, cfg); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
