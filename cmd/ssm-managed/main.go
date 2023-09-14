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

package main

import (
	"bytes"
	"database/sql"
	_ "expvar"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/AlekSi/pointer"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/revel/config"
	pc "github.com/shatteredsilicon/ssm/proto/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/reform.v1"
	reformMySQL "gopkg.in/reform.v1/dialects/mysql"

	"github.com/shatteredsilicon/ssm-managed/api"
	"github.com/shatteredsilicon/ssm-managed/handlers"
	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/services/consul"
	"github.com/shatteredsilicon/ssm-managed/services/grafana"
	"github.com/shatteredsilicon/ssm-managed/services/logs"
	"github.com/shatteredsilicon/ssm-managed/services/mysql"
	"github.com/shatteredsilicon/ssm-managed/services/node"
	"github.com/shatteredsilicon/ssm-managed/services/postgresql"
	"github.com/shatteredsilicon/ssm-managed/services/prometheus"
	"github.com/shatteredsilicon/ssm-managed/services/qan"
	"github.com/shatteredsilicon/ssm-managed/services/rds"
	"github.com/shatteredsilicon/ssm-managed/services/remote"
	"github.com/shatteredsilicon/ssm-managed/services/snmp"
	"github.com/shatteredsilicon/ssm-managed/services/supervisor"
	"github.com/shatteredsilicon/ssm-managed/services/telemetry"
	"github.com/shatteredsilicon/ssm-managed/utils"
	"github.com/shatteredsilicon/ssm-managed/utils/interceptors"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
	"github.com/shatteredsilicon/ssm-managed/utils/ports"
)

const (
	shutdownTimeout = 3 * time.Second
)

var (
	// TODO we can combine gRPC and REST ports, but only with TLS
	// see https://github.com/grpc/grpc-go/issues/555
	// alternatively, we can try to use cmux: https://open.dgraph.io/post/cmux/
	gRPCAddrF  = flag.String("listen-grpc-addr", "127.0.0.1:7771", "gRPC server listen address")
	restAddrF  = flag.String("listen-rest-addr", "127.0.0.1:7772", "REST server listen address")
	debugAddrF = flag.String("listen-debug-addr", "127.0.0.1:7773", "Debug server listen address")

	swaggerF = flag.String("swagger", "off", "Server to serve Swagger: rest, debug or off")

	prometheusConfigF = flag.String("prometheus-config", "", "Prometheus configuration file path")
	prometheusURLF    = flag.String("prometheus-url", "http://127.0.0.1:9090/", "Prometheus base URL")
	promtoolF         = flag.String("promtool", "promtool", "promtool path")

	consulAddrF  = flag.String("consul-addr", "127.0.0.1:8500", "Consul HTTP API address")
	grafanaAddrF = flag.String("grafana-addr", "127.0.0.1:3000", "Grafana HTTP API address")

	dbNameF       = flag.String("db-name", "", "Database name")
	dbUsernameF   = flag.String("db-username", "ssm-managed", "Database username")
	dbPasswordF   = flag.String("db-password", "ssm-managed", "Database password")
	dbSocketF     = flag.String("db-socket", "/var/lib/mysql/mysql.sock", "Database socket")
	qanAPIConfigF = flag.String("qan-api-config", "/etc/ssm-qan-api.conf", "QAN configuration file path")

	agentMySQLdExporterF    = flag.String("agent-mysqld-exporter", "/opt/ss/ssm-client/mysqld_exporter", "mysqld_exporter path")
	agentPostgresExporterF  = flag.String("agent-postgres-exporter", "/opt/ss/ssm-client/postgres_exporter", "postgres_exporter path")
	agentRDSExporterF       = flag.String("agent-rds-exporter", "/usr/sbin/rds_exporter", "rds_exporter path")
	agentRDSExporterConfigF = flag.String("agent-rds-exporter-config", "/etc/ssm-rds-exporter.yml", "rds_exporter configuration file path")
	agentSNMPExporterF      = flag.String("agent-snmp-exporter", "/opt/ss/snmp_exporter/bin/snmp_exporter", "snmp_exporter path")
	snmpGeneratorF          = flag.String("snmp-generator", "/opt/ss/snmp_exporter/bin/generator", "snmp_generator path")
	agentSNMPConfigDirF     = flag.String("agent-snmp-config-dir", "/opt/ss/snmp_exporter", "snmp_exporter config directory")
	agentQANBaseF           = flag.String("agent-qan-base", "/opt/ss/qan-agent", "qan-agent installation base path")

	rdsEnableGovCloud = flag.Bool("rds-enable-gov-cloud", false, "Enable GOV cloud for RDS")
	rdsEnableCnCloud  = flag.Bool("rds-enable-cn-cloud", false, "Enable AWS CN cloud for RDS")

	debugF = flag.Bool("debug", false, "Enable debug logging")
)

func addSwaggerHandler(mux *http.ServeMux) {
	// TODO embed swagger resources?
	pattern := "/swagger/"
	fileServer := http.StripPrefix(pattern, http.FileServer(http.Dir("api/swagger")))
	mux.HandleFunc(pattern, func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		fileServer.ServeHTTP(rw, req)
	})
}

func addLogsHandler(mux *http.ServeMux, logs *logs.Logs) {
	l := logrus.WithField("component", "logs.zip")

	mux.HandleFunc("/logs.zip", func(rw http.ResponseWriter, req *http.Request) {
		// fail-safe
		ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
		defer cancel()

		t := time.Now().UTC()
		filename := fmt.Sprintf("ssm-server_%4d-%02d-%02d-%02d-%02d.zip", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())

		rw.Header().Set(`Access-Control-Allow-Origin`, `*`)
		rw.Header().Set(`Content-Type`, `application/zip`)
		rw.Header().Set(`Content-Disposition`, `attachment; filename="`+filename+`"`)
		ctx, _ = logger.Set(ctx, "logs")
		if err := logs.Zip(ctx, rw); err != nil {
			l.Error(err)
		}
	})
}

func makePortsRegistry(db *reform.DB) (*ports.Registry, error) {
	// collect already reserved ports
	rows, err := db.Query("SELECT listen_port FROM agents")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var reserved []uint16
	for rows.Next() {
		var port uint16
		if err = rows.Scan(&port); err != nil {
			return nil, errors.WithStack(err)
		}
		reserved = append(reserved, port)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	registry := ports.NewRegistry(10000, 10999, reserved)
	return registry, err
}

func makeInternalQan(ctx context.Context, deps *serviceDependencies) error {
	var nodeID int32

	agentService, err := models.AgentServiceByName(deps.db.Querier, models.SSMServerNodeName, string(models.QanAgentAgentType))
	if err != nil {
		return err
	}
	if agentService == nil {
		var node models.Node
		err := deps.db.FindOneTo(&node, "type", models.SSMServerNodeType)
		if err != nil {
			return err
		}

		nodeID = node.ID
	} else {
		nodeID = agentService.NodeID
	}

	var version, versionComment string
	err = deps.db.QueryRowContext(ctx, "SELECT @@version, @@version_comment").Scan(&version, &versionComment)
	if err != nil {
		return err
	}

	engine, engineVersion, err := mysql.NormalizeEngineAndEngineVersion(versionComment, version)
	if err != nil {
		return err
	}

	service := models.MySQLService{
		Type:          models.MySQLServiceType,
		NodeID:        nodeID,
		Address:       dbSocketF,
		Engine:        &engine,
		EngineVersion: &engineVersion,
	}

	serviceUsername := "root"
	servicePassword := ""
	agent := models.QanAgent{
		Type:            models.QanAgentAgentType,
		RunsOnNodeID:    nodeID,
		ServiceUsername: &serviceUsername,
		ServicePassword: &servicePassword,
		ListenPort:      pointer.ToUint16(models.QanAgentPort),
	}

	err = deps.qan.AddMySQL(
		ctx,
		string(models.SSMServerNodeType),
		&service,
		&agent,
		pc.QAN{
			CollectFrom: qan.SlowlogCollectFrom,
			FilterAllow: []string{"SELECT", "DELETE"},
		},
	)
	if err != nil {
		return err
	}

	if agentService != nil {
		return nil
	}

	return deps.db.InTransaction(func(tx *reform.TX) error {
		if err = tx.Insert(&service); err != nil {
			return err
		}

		if err = tx.Insert(&agent); err != nil {
			return err
		}

		return tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID})
	})
}

func removeInternalQan(ctx context.Context, deps *serviceDependencies) error {
	agentService, err := models.AgentServiceByName(deps.db.Querier, models.SSMServerNodeName, string(models.QanAgentAgentType))
	if err != nil {
		return err
	}

	if agentService == nil {
		return nil
	}

	agent := models.QanAgent{
		Type:              models.QanAgentAgentType,
		RunsOnNodeID:      agentService.NodeID,
		ListenPort:        pointer.ToUint16(models.QanAgentPort),
		QANDBInstanceUUID: &agentService.QanDBInstanceUUID,
	}

	err = deps.qan.RemoveMySQL(ctx, &agent)
	if err != nil {
		return err
	}

	err = deps.qan.RemoveQANData(ctx, agentService.QanDBInstanceUUID)
	if err != nil {
		return err
	}

	return deps.db.InTransaction(func(tx *reform.TX) error {
		_, err = tx.DeleteFrom(models.AgentServiceView, "WHERE agent_id = ?", agentService.AgentID)
		if err != nil {
			return err
		}

		_, err = tx.DeleteFrom(models.AgentTable, "WHERE id = ?", agentService.AgentID)
		if err != nil {
			return err
		}

		_, err = tx.DeleteFrom(models.ServiceTable, "WHERE id = ?", agentService.ServiceID)
		return err
	})
}

type serviceDependencies struct {
	prometheus    *prometheus.Service
	supervisor    *supervisor.Supervisor
	db            *reform.DB
	portsRegistry *ports.Registry
	qan           *qan.Service
}

func makeRDSService(ctx context.Context, deps *serviceDependencies) (*rds.Service, error) {
	rdsConfig := rds.ServiceConfig{
		MySQLdExporterPath:    *agentMySQLdExporterF,
		RDSExporterPath:       *agentRDSExporterF,
		RDSExporterConfigPath: *agentRDSExporterConfigF,

		Prometheus:    deps.prometheus,
		Supervisor:    deps.supervisor,
		DB:            deps.db,
		PortsRegistry: deps.portsRegistry,
		QAN:           deps.qan,

		RDSEnableGovCloud: *rdsEnableGovCloud,
		RDSEnableCnCloud:  *rdsEnableCnCloud,
	}
	rdsService, err := rds.NewService(&rdsConfig)
	if err != nil {
		return nil, err
	}

	err = deps.db.InTransaction(func(tx *reform.TX) error {
		return rdsService.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})
	if err != nil {
		return nil, err
	}
	err = deps.db.InTransaction(func(tx *reform.TX) error {
		return rdsService.Restore(ctx, tx)
	})
	if err != nil {
		return nil, err
	}

	return rdsService, nil
}

func makeSNMPService(ctx context.Context, deps *serviceDependencies, consul *consul.Client) (*snmp.Service, error) {
	snmpConfig := snmp.ServiceConfig{
		SNMPExporterPath:  *agentSNMPExporterF,
		SNMPGeneratorPath: *snmpGeneratorF,
		SNMPConfigDir:     *agentSNMPConfigDirF,

		Prometheus:    deps.prometheus,
		Supervisor:    deps.supervisor,
		DB:            deps.db,
		PortsRegistry: deps.portsRegistry,
		Consul:        consul,
	}
	snmpService, err := snmp.NewService(&snmpConfig)
	if err != nil {
		return nil, err
	}

	err = deps.db.InTransaction(func(tx *reform.TX) error {
		return snmpService.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})
	if err != nil {
		return nil, err
	}

	return snmpService, nil
}

func makeMySQLService(ctx context.Context, deps *serviceDependencies, consul *consul.Client) (*mysql.Service, error) {
	serviceConfig := mysql.ServiceConfig{
		MySQLdExporterPath: *agentMySQLdExporterF,

		Prometheus:    deps.prometheus,
		Supervisor:    deps.supervisor,
		DB:            deps.db,
		PortsRegistry: deps.portsRegistry,
		QAN:           deps.qan,
		Consul:        consul,
	}
	mysqlService, err := mysql.NewService(&serviceConfig)
	if err != nil {
		return nil, err
	}

	err = deps.db.InTransaction(func(tx *reform.TX) error {
		return mysqlService.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})
	if err != nil {
		return nil, err
	}
	err = deps.db.InTransaction(func(tx *reform.TX) error {
		return mysqlService.Restore(ctx, tx)
	})
	if err != nil {
		return nil, err
	}

	return mysqlService, nil
}

func makePostgreSQLService(ctx context.Context, deps *serviceDependencies, consul *consul.Client) (*postgresql.Service, error) {
	serviceConfig := postgresql.ServiceConfig{
		PostgresExporterPath: *agentPostgresExporterF,

		Prometheus:    deps.prometheus,
		Supervisor:    deps.supervisor,
		DB:            deps.db,
		PortsRegistry: deps.portsRegistry,
		Consul:        consul,
	}
	postgresqlService, err := postgresql.NewService(&serviceConfig)
	if err != nil {
		return nil, err
	}

	err = deps.db.InTransaction(func(tx *reform.TX) error {
		return postgresqlService.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})
	if err != nil {
		return nil, err
	}
	err = deps.db.InTransaction(func(tx *reform.TX) error {
		return postgresqlService.Restore(ctx, tx)
	})
	if err != nil {
		return nil, err
	}

	return postgresqlService, nil
}

type grpcServerDependencies struct {
	*serviceDependencies
	consulClient *consul.Client
	rds          *rds.Service
	mysql        *mysql.Service
	postgres     *postgresql.Service
	snmp         *snmp.Service
	remote       *remote.Service
	logs         *logs.Logs
	node         *node.Service
}

// runGRPCServer runs gRPC server until context is canceled, then gracefully stops it.
func runGRPCServer(ctx context.Context, deps *grpcServerDependencies) {
	l := logrus.WithField("component", "gRPC")
	l.Infof("Starting server on http://%s/ ...", *gRPCAddrF)

	grafana := grafana.NewClient(*grafanaAddrF)

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.Unary),
		grpc.StreamInterceptor(interceptors.Stream),
	)
	api.RegisterBaseServer(gRPCServer, &handlers.BaseServer{PMMVersion: utils.Version})
	api.RegisterDemoServer(gRPCServer, &handlers.DemoServer{})
	api.RegisterScrapeConfigsServer(gRPCServer, &handlers.ScrapeConfigsServer{
		Prometheus: deps.prometheus,
	})
	api.RegisterRDSServer(gRPCServer, &handlers.RDSServer{
		RDS: deps.rds,
	})
	api.RegisterMySQLServer(gRPCServer, &handlers.MySQLServer{
		MySQL: deps.mysql,
	})
	api.RegisterPostgreSQLServer(gRPCServer, &handlers.PostgreSQLServer{
		PostgreSQL: deps.postgres,
	})
	api.RegisterSNMPServer(gRPCServer, &handlers.SNMPServer{
		SNMP: deps.snmp,
	})
	api.RegisterRemoteServer(gRPCServer, &handlers.RemoteServer{
		Remote: deps.remote,
	})
	api.RegisterLogsServer(gRPCServer, &handlers.LogsServer{
		Logs: deps.logs,
	})
	api.RegisterAnnotationsServer(gRPCServer, &handlers.AnnotationsServer{
		Grafana: grafana,
	})
	api.RegisterNodeServer(gRPCServer, &handlers.NodeServer{
		Node:   deps.node,
		Remote: deps.remote,
	})

	grpc_prometheus.Register(gRPCServer)
	grpc_prometheus.EnableHandlingTimeHistogram()

	listener, err := net.Listen("tcp", *gRPCAddrF)
	if err != nil {
		l.Panic(err)
	}
	go func() {
		for {
			err = gRPCServer.Serve(listener)
			if err == nil || err == grpc.ErrServerStopped {
				break
			}
			l.Errorf("Failed to serve: %s", err)
		}
		l.Info("Server stopped.")
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	go func() {
		<-ctx.Done()
		gRPCServer.Stop()
	}()
	gRPCServer.GracefulStop()
	cancel()
}

// runRESTServer runs REST proxy server until context is canceled, then gracefully stops it.
func runRESTServer(ctx context.Context, logs *logs.Logs) {
	l := logrus.WithField("component", "REST")
	l.Infof("Starting server on http://%s/ ...", *restAddrF)

	proxyMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	type registrar func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
	for _, r := range []registrar{
		api.RegisterBaseHandlerFromEndpoint,
		api.RegisterDemoHandlerFromEndpoint,
		api.RegisterScrapeConfigsHandlerFromEndpoint,
		api.RegisterRDSHandlerFromEndpoint,
		api.RegisterMySQLHandlerFromEndpoint,
		api.RegisterPostgreSQLHandlerFromEndpoint,
		api.RegisterSNMPHandlerFromEndpoint,
		api.RegisterRemoteHandlerFromEndpoint,
		api.RegisterLogsHandlerFromEndpoint,
		api.RegisterAnnotationsHandlerFromEndpoint,
		api.RegisterNodeHandlerFromEndpoint,
	} {
		if err := r(ctx, proxyMux, *gRPCAddrF, opts); err != nil {
			l.Panic(err)
		}
	}

	mux := http.NewServeMux()
	if *swaggerF == "rest" {
		l.Printf("Swagger enabled. http://%s/swagger/", *restAddrF)
		addSwaggerHandler(mux)
	}
	addLogsHandler(mux, logs)
	mux.Handle("/", proxyMux)

	server := &http.Server{
		Addr:     *restAddrF,
		ErrorLog: log.New(os.Stderr, "runRESTServer: ", 0),
		Handler:  mux,

		// TODO we probably will need it for TLS+HTTP/2, see https://github.com/philips/grpc-gateway-example/issues/11
		// TLSConfig: &tls.Config{
		// 	NextProtos: []string{"h2"},
		// },
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			l.Panic(err)
		}
		l.Info("Server stopped.")
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	if err := server.Shutdown(ctx); err != nil {
		l.Errorf("Failed to shutdown gracefully: %s", err)
	}
	cancel()
}

// runDebugServer runs debug server until context is canceled, then gracefully stops it.
func runDebugServer(ctx context.Context) {
	l := logrus.WithField("component", "debug")

	http.Handle("/debug/metrics", promhttp.Handler())

	handlers := []string{"/debug/metrics", "/debug/vars", "/debug/requests", "/debug/events", "/debug/pprof"}
	if *swaggerF == "debug" {
		handlers = append(handlers, "/swagger")
		l.Printf("Swagger enabled. http://%s/swagger/", *debugAddrF)
		addSwaggerHandler(http.DefaultServeMux)
	}

	for i, h := range handlers {
		handlers[i] = "http://" + *debugAddrF + h
	}

	var buf bytes.Buffer
	err := template.Must(template.New("debug").Parse(`
	<html>
	<body>
	<ul>
	{{ range . }}
		<li><a href="{{ . }}">{{ . }}</a></li>
	{{ end }}
	</ul>
	</body>
	</html>
	`)).Execute(&buf, handlers)
	if err != nil {
		l.Panic(err)
	}
	http.HandleFunc("/debug", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(buf.Bytes())
	})
	l.Infof("Starting server on http://%s/debug\nRegistered handlers:\n\t%s", *debugAddrF, strings.Join(handlers, "\n\t"))

	server := &http.Server{
		Addr:     *debugAddrF,
		ErrorLog: log.New(os.Stderr, "runDebugServer: ", 0),
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			l.Panic(err)
		}
		l.Info("Server stopped.")
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	if err := server.Shutdown(ctx); err != nil {
		l.Errorf("Failed to shutdown gracefully: %s", err)
	}
	cancel()
}

func runTelemetryService(ctx context.Context, consulClient *consul.Client) {
	l := logrus.WithField("component", "telemetry")

	uuid, err := getTelemetryUUID(consulClient)
	if err != nil {
		l.Panicf("cannot get/set telemetry UUID in Consul: %s", err)
	}

	svc := telemetry.NewService(uuid, utils.Version)
	svc.Run(ctx)
}

func getTelemetryUUID(consulClient *consul.Client) (string, error) {
	b, err := consulClient.GetKV("telemetry/uuid")
	if err != nil {
		return "", err
	}
	if len(b) > 0 {
		return string(b), nil
	}

	uuid, err := telemetry.GenerateUUID()
	if err != nil {
		return "", err
	}
	if err = consulClient.PutKV("telemetry/uuid", []byte(uuid)); err != nil {
		return "", err
	}
	return uuid, nil
}

func main() {
	log.SetFlags(0)
	log.Printf("ssm-managed %s", utils.Version)
	log.SetPrefix("stdlog: ")
	flag.Parse()

	if *dbNameF == "" {
		log.Fatal("-db-name flag must be given explicitly.")
	}

	if *debugF {
		logrus.SetLevel(logrus.DebugLevel)
		grpclog.SetLoggerV2(&logger.GRPC{Entry: logrus.WithField("component", "grpclog")})
	}

	if *swaggerF != "rest" && *swaggerF != "debug" && *swaggerF != "off" {
		flag.Usage()
		log.Fatalf("Unexpected value %q for -swagger flag.", *swaggerF)
	}

	l := logrus.WithField("component", "main")
	ctx, cancel := context.WithCancel(context.Background())
	ctx, _ = logger.Set(ctx, "main")
	defer l.Info("Done.")

	// handle termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-signals
		signal.Stop(signals)
		l.Warnf("Got %v (%d) signal, shutting down...", s, s)
		cancel()
	}()

	consulClient, err := consul.NewClient(*consulAddrF)
	if err != nil {
		l.Panic(err)
	}

	prometheus, err := prometheus.NewService(*prometheusConfigF, *prometheusURLF, *promtoolF, consulClient)
	if err == nil {
		err = prometheus.Check(ctx)
	}
	if err != nil {
		l.Panicf("Prometheus service problem: %+v", err)
	}

	supervisor := supervisor.New(l)

	// open QAN db conn
	qanConfig, err := config.ReadDefault(*qanAPIConfigF)
	if err != nil {
		l.Panic(err)
	}
	qanDSN, err := qanConfig.RawStringDefault("mysql.dsn")
	if err != nil {
		l.Panic(err)
	}
	qanDB, err := sql.Open("mysql", qanDSN)
	if err != nil {
		l.Panic(err)
	}

	qan, err := qan.NewService(ctx, *agentQANBaseF, supervisor, qanDB)
	if err != nil {
		l.Panicf("QAN service problem: %+v", err)
	}

	sqlDB, err := models.OpenDB(*dbNameF, *dbUsernameF, *dbPasswordF, l.Debugf)
	if err != nil {
		l.Panic(err)
	}
	defer sqlDB.Close()
	db := reform.NewDB(sqlDB, reformMySQL.Dialect, nil)

	portsRegistry, err := makePortsRegistry(db)
	if err != nil {
		l.Panic(err)
	}

	deps := &serviceDependencies{
		prometheus:    prometheus,
		supervisor:    supervisor,
		qan:           qan,
		db:            db,
		portsRegistry: portsRegistry,
	}

	if os.Getenv("DEBUG") == "1" || os.Getenv("DEBUG") == "true" {
		err = makeInternalQan(ctx, deps)
		if err != nil {
			l.Panicf("Failed to add internal qan: %+v", err)
		}
	} else {
		err = removeInternalQan(ctx, deps)
		if err != nil {
			l.Warnf("Failed to remove internal qan: %+v", err)
		}
	}

	// restore all qan configs from database before
	// restoring rds or remote mysql services
	err = deps.qan.RestoreConfigs(ctx, deps.db.Querier)
	if err != nil {
		l.Panicf("Restore qan configs failed: %+v", err)
	}

	rds, err := makeRDSService(ctx, deps)
	if err != nil {
		l.Panicf("RDS service problem: %+v", err)
	}

	mysqlService, err := makeMySQLService(ctx, deps, consulClient)
	if err != nil {
		l.Panicf("MySQL service problem: %+v", err)
	}

	postgres, err := makePostgreSQLService(ctx, deps, consulClient)
	if err != nil {
		l.Panicf("PostgreSQL service problem: %+v", err)
	}

	snmp, err := makeSNMPService(ctx, deps, consulClient)
	if err != nil {
		l.Panicf("SNMP service problem: %+v", err)
	}

	remoteService, err := remote.NewService(&remote.ServiceConfig{
		DB: deps.db,
	})
	if err != nil {
		l.Panicf("Remote service problem: %+v", err)
	}

	logs := logs.New(utils.Version, consulClient, db, rds, nil)

	nodeService := node.NewService(consulClient, deps.qan, deps.prometheus, deps.db, mysqlService, postgres, rds, snmp)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		runGRPCServer(ctx, &grpcServerDependencies{
			serviceDependencies: deps,
			rds:                 rds,
			postgres:            postgres,
			mysql:               mysqlService,
			snmp:                snmp,
			remote:              remoteService,
			consulClient:        consulClient,
			logs:                logs,
			node:                nodeService,
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runRESTServer(ctx, logs)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runDebugServer(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runTelemetryService(ctx, consulClient)
	}()

	wg.Wait()
}
