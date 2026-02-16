package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"sort"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/go-sql-driver/mysql"
	servicelib "github.com/percona/kardianos-service"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/services"
	"github.com/shatteredsilicon/ssm-managed/services/consul"
	"github.com/shatteredsilicon/ssm-managed/services/prometheus"
	"github.com/shatteredsilicon/ssm-managed/services/qan"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
	"github.com/shatteredsilicon/ssm-managed/utils/ports"
	"github.com/shatteredsilicon/ssm/proto/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var defaultRemoteQANConfig = config.QAN{}

type ServiceConfig struct {
	MongoDBExporterPath string

	Prometheus    *prometheus.Service
	Supervisor    services.Supervisor
	DB            *reform.DB
	QAN           *qan.Service
	PortsRegistry *ports.Registry
	Consul        *consul.Client
}

// Service is responsible for interactions with MongoDB.
type Service struct {
	*ServiceConfig
	ssmServerNode *models.Node
}

// NewService creates a new service.
func NewService(config *ServiceConfig) (*Service, error) {
	var node models.Node
	err := config.DB.FindOneTo(&node, "type", models.SSMServerNodeType)
	if err != nil {
		return nil, err
	}

	for _, path := range []*string{
		&config.MongoDBExporterPath,
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
		ssmServerNode: &node,
	}
	return svc, nil
}

// ApplyPrometheusConfiguration Adds mongodb to prometheus configuration and applies it
func (svc *Service) ApplyPrometheusConfiguration(ctx context.Context, q *reform.Querier) error {
	mongodbConfig := &prometheus.ScrapeConfig{
		JobName:        "remote-mongodb",
		ScrapeInterval: "1s",
		ScrapeTimeout:  "1s",
		MetricsPath:    "/metrics",
		RelabelConfigs: []prometheus.RelabelConfig{{
			TargetLabel: "job",
			Replacement: "mongodb",
		}},
		MetricRelabelConfigs: []prometheus.RelabelConfig{
			{
				SourceLabels: []model.LabelName{"__name__"},
				TargetLabel:  "__name__",
				Regex:        "(mongodb_mongod_durability_journaled_megabytes|mongodb_mongod_durability_write_to_data_files_megabytes|mongodb_mongod_durability_commits|mongodb_mongod_background_flushing_average_milliseconds|mongodb_mongod_metrics_get_last_error_wtime_total_milliseconds|mongodb_mongos_metrics_get_last_error_wtime_total_milliseconds|mongodb_mongod_metrics_repl_network_getmores_total_milliseconds|mongodb_mongod_metrics_repl_preload_docs_total_milliseconds|mongodb_mongod_metrics_repl_preload_indexes_total_milliseconds|mongodb_mongod_metrics_repl_apply_batches_total_milliseconds|mongodb_mongod_rocksdb_compaction_bytes_per_second)",
				Replacement:  "${1}_total",
			},
		},
	}

	nodes, err := q.FindAllFrom(models.RemoteNodeTable, "type", models.RemoteNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.RemoteNode)

		mongodbServices, e := q.SelectAllFrom(models.MongoDBServiceTable, "WHERE node_id = ? AND type = ?", node.ID, models.MongoDBServiceType)
		if e != nil {
			return errors.WithStack(e)
		}
		if len(mongodbServices) == 0 {
			continue
		}
		service := mongodbServices[0].(*models.MongoDBService)
		if service.Type != models.MongoDBServiceType {
			continue
		}

		agents, err := models.AgentsForServiceID(q, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.MongoDBExporterAgentType:
				a := models.MongoDBExporter{ID: agent.ID}
				if e := q.Reload(&a); e != nil {
					return errors.WithStack(e)
				}
				logger.Get(ctx).WithField("component", "mongodb").Infof("%s %s %d", a.Type, node.Name, *a.ListenPort)

				sc := prometheus.StaticConfig{
					Targets: []string{fmt.Sprintf("127.0.0.1:%d", *a.ListenPort)},
					Labels: []prometheus.LabelPair{
						{Name: "instance", Value: node.Name},
						{Name: "region", Value: string(models.RemoteNodeRegion)},
					},
				}
				mongodbConfig.StaticConfigs = append(mongodbConfig.StaticConfigs, sc)
			}
		}
	}

	// sort by instance
	sorterFor := func(sc []prometheus.StaticConfig) func(int, int) bool {
		return func(i, j int) bool {
			return sc[i].Labels[0].Value < sc[j].Labels[0].Value
		}
	}
	sort.Slice(mongodbConfig.StaticConfigs, sorterFor(mongodbConfig.StaticConfigs))

	return svc.Prometheus.SetScrapeConfigs(ctx, false, mongodbConfig)
}

// Add new MongoDB service and start mongodb_exporter
func (svc *Service) Add(
	ctx context.Context,
	name, uri string,
	qanConfig *config.QAN,
) (int32, error) {
	uri = strings.TrimSpace(uri)
	name = strings.TrimSpace(name)
	if uri == "" {
		return 0, status.Error(codes.InvalidArgument, "MongoDB uri is not given.")
	}

	clientOpts, err := GetClientOpts(uri)
	if err != nil {
		return 0, err
	}

	if name == "" {
		name = strings.Join(clientOpts.Hosts, ",")
	}

	// check if instance is added on client side
	added, err := svc.clientInstanceAdded(ctx, name)
	if err != nil {
		return 0, err
	}
	if added {
		return 0, status.Error(codes.AlreadyExists, fmt.Sprintf("MongoDB instance with name %s is already added", name))
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

		engine, engineVersion, err := svc.EngineAndEngineVersion(ctx, clientOpts)
		if err != nil {
			return errors.WithStack(err)
		}

		// insert service
		service := &models.MongoDBService{
			Type:   models.MongoDBServiceType,
			NodeID: node.ID,

			Address:       &uri,
			Engine:        &engine,
			EngineVersion: &engineVersion,
		}
		if err := tx.Insert(service); err != nil {
			return errors.WithStack(err)
		}

		if err := svc.AddMongoDBExporter(ctx, tx, service, clientOpts); err != nil {
			return err
		}
		if err = svc.addQanAgent(ctx, tx, service, node, uri, qanConfig); err != nil {
			return err
		}

		return svc.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})

	return id, err
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
		if t == models.ClientMongoDBExporterAgentType {
			// instance is added on client side
			return true, nil
		}
	}

	return false, nil
}

// BuildInfo response struct of mongo buildInfo command
type BuildInfo struct {
	Version        string `bson:"version"`
	VersionArray   []int  `bson:"versionArray"`
	GitVersion     string `bson:"gitVersion"`
	OpenSSLVersion string `bson:"OpenSSLVersion"`
	SysInfo        string `bson:"sysInfo"`
	Bits           int    `bson:"bits"`
	Debug          bool   `bson:"debug"`
	MaxObjectSize  int    `bson:"maxBsonObjectSize"`
}

// EngineAndEngineVersion returns mongodb engine and version
func (svc *Service) EngineAndEngineVersion(ctx context.Context, opts *options.ClientOptions) (string, string, error) {
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return "", "", err
	}
	defer client.Disconnect(ctx)

	var info BuildInfo
	err = client.Database("admin").RunCommand(ctx, bson.D{{"buildInfo", "1"}}).Decode(&info)
	if err != nil {
		return "", "", err
	}

	return "MongoDB", info.Version, nil
}

func (svc *Service) AddMongoDBExporter(ctx context.Context, tx *reform.TX, service *models.MongoDBService, mongoOpts *options.ClientOptions) error {
	// insert mongodb_exporter agent and association
	port, err := svc.PortsRegistry.Reserve()
	if err != nil {
		return err
	}
	agent := &models.MongoDBExporter{
		Type:         models.MongoDBExporterAgentType,
		RunsOnNodeID: svc.ssmServerNode.ID,

		ListenPort: &port,
	}
	if err = tx.Insert(agent); err != nil {
		return errors.WithStack(err)
	}
	if err = tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID}); err != nil {
		return errors.WithStack(err)
	}

	// check connection
	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		return errors.WithStack(err)
	}
	defer client.Disconnect(ctx)

	if err = client.Ping(ctx, nil); err != nil {
		return errors.WithStack(err)
	}

	// start mongodb_exporter agent
	if svc.MongoDBExporterPath != "" {
		cfg := svc.MongoDBExporterCfg(agent, mongoOpts.GetURI())
		if err = svc.Supervisor.Start(ctx, cfg); err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) addQanAgent(
	ctx context.Context,
	tx *reform.TX,
	service *models.MongoDBService,
	node *models.RemoteNode,
	uri string,
	qanConfig *config.QAN,
) error {
	// Despite running a single qan-agent process on PMM Server, we use one database record per MySQL instance
	// to store username/password and UUID.

	// insert qan-agent agent and association
	agent := &models.QanAgent{
		Type:         models.QanAgentAgentType,
		RunsOnNodeID: svc.ssmServerNode.ID,

		ListenPort: pointer.ToUint16(models.QanAgentPort),
	}
	var err error
	if err = tx.Insert(agent); err != nil {
		return errors.WithStack(err)
	}
	if err = tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID}); err != nil {
		return errors.WithStack(err)
	}

	// DSNs for mysqld_exporter and qan-agent are currently identical,
	// so we do not check connection again

	// start or reconfigure qan-agent
	if svc.QAN != nil {
		if qanConfig == nil {
			qanConfig = &defaultRemoteQANConfig
		}

		nodeName := node.Name
		if node.Type == models.SSMServerNodeType {
			nodeName = string(node.Type) // ssm-server node uses type as name
		}
		if err = svc.QAN.AddQAN(ctx, nodeName, "mongo", uri, *service.EngineVersion, agent, *qanConfig); err != nil {
			return err
		}

		// re-save agent with set QANDBInstanceUUID
		if err = tx.Save(agent); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// Restore configuration from database.
func (svc *Service) Restore(ctx context.Context, tx *reform.TX) error {
	nodes, err := tx.FindAllFrom(models.RemoteNodeTable, "type", models.RemoteNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.RemoteNode)

		mongodbServices, e := tx.SelectAllFrom(models.MongoDBServiceTable, "WHERE node_id = ? AND type = ?", node.ID, models.MongoDBServiceType)
		if e != nil {
			return errors.WithStack(e)
		}
		if len(mongodbServices) == 0 {
			continue
		}

		service := mongodbServices[0].(*models.MongoDBService)
		agents, err := models.AgentsForServiceID(tx.Querier, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.MongoDBExporterAgentType:
				a := &models.MongoDBExporter{ID: agent.ID}
				if err = tx.Reload(a); err != nil {
					return errors.WithStack(err)
				}
				if svc.MongoDBExporterPath != "" {
					name := models.NameForSupervisor(a.Type, *a.ListenPort)

					err := svc.Supervisor.Status(ctx, name)
					if err == nil {
						if err = svc.Supervisor.Stop(ctx, name); err != nil {
							return err
						}
					}

					cfg := svc.MongoDBExporterCfg(a, *service.Address)
					if err = svc.Supervisor.Start(ctx, cfg); err != nil {
						return err
					}
				}

			case models.QanAgentAgentType:
				a := models.QanAgent{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.QAN != nil {
					name := models.NameForSupervisor(a.Type, *a.ListenPort)
					err := svc.Supervisor.Status(ctx, name)
					if err == nil {
						if err = svc.Supervisor.Stop(ctx, name); err != nil {
							return err
						}
					}

					if err = svc.QAN.Restore(ctx, name, models.QanAgentWithSubsystem{QanAgent: a, Subsystem: "mongo"}, defaultRemoteQANConfig); err != nil {
						if _, ok := err.(qan.QANCommandError); ok {
							// if it's a QAN command error, we should have already
							// restored the qan configs (although may not be perfectly),
							// one should check what happens on the qan-agent side, ssm-managed
							// should just continue on.
							logrus.WithField("component", "rds").Warnf("Got a QAN API error when restoring qan for %s: %s\n", node.Name, err.Error())
							return nil
						}

						return err
					}
				}
			}
		}
	}

	return nil
}

func (svc *Service) MongoDBExporterCfg(agent *models.MongoDBExporter, uri string) *servicelib.Config {
	name := models.NameForSupervisor(agent.Type, *agent.ListenPort)

	arguments := []string{
		fmt.Sprintf("-web.listen-address=127.0.0.1:%d", *agent.ListenPort),
		"-web.auth-file=\"\"",
	}
	sort.Strings(arguments)

	return &servicelib.Config{
		Name:        name,
		DisplayName: name,
		Description: name,
		Executable:  svc.MongoDBExporterPath,
		Arguments:   arguments,
		Environment: []string{fmt.Sprintf("MONGODB_URI=%s", uri)},
	}
}

func GetClientOpts(uri string) (*options.ClientOptions, error) {
	connStr, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, err
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// default to direct connect when possible
	if len(opts.Hosts) == 1 && connStr.Scheme == connstring.SchemeMongoDB && opts.Direct == nil {
		opts.SetDirect(true)
	}

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return opts, nil
}

func SanitizeURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return uri
	}

	u.User = nil
	return u.String()
}
