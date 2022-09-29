package node

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/services"
	"github.com/shatteredsilicon/ssm-managed/services/consul"
	"github.com/shatteredsilicon/ssm-managed/services/mysql"
	"github.com/shatteredsilicon/ssm-managed/services/postgresql"
	"github.com/shatteredsilicon/ssm-managed/services/prometheus"
	"github.com/shatteredsilicon/ssm-managed/services/qan"
	"github.com/shatteredsilicon/ssm-managed/services/rds"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

const (
	removedInstanceConsulKey = "client/removed-instances"
)

// Service client service
type Service struct {
	consul     *consul.Client
	qan        *qan.Service
	prometheus *prometheus.Service
	db         *reform.DB
	mysql      *mysql.Service
	postgresql *postgresql.Service
	rds        *rds.Service
}

func NewService(
	consul *consul.Client, qan *qan.Service, prometheus *prometheus.Service,
	db *reform.DB, mysql *mysql.Service, postgresql *postgresql.Service,
	rds *rds.Service,
) *Service {
	return &Service{
		consul:     consul,
		qan:        qan,
		prometheus: prometheus,
		db:         db,
		mysql:      mysql,
		postgresql: postgresql,
		rds:        rds,
	}
}

// RemovedInstance removed instance structure
type RemovedInstance struct {
	Name string
}

// ClientNodeService service of client-added node
type ClientNodeService struct {
	Name    string
	Address string
	Port    int
	Region  string
	Distro  string
	Version string
}

// ClientNode client-added node
type ClientNode struct {
	Name     string
	Services []ClientNodeService
}

// RemoveNode removes node
func (svc *Service) RemoveNode(ctx context.Context, nodeID string) error {
	// remove qan data
	err := svc.removeNodeFromQan(ctx, nodeID)
	if err != nil {
		logrus.Errorf("delete qan data for %s failed: %+v", nodeID, err)
		return err
	}

	// remove consul data, for client-added nodes
	err = svc.removeNodeFromConsul(ctx, nodeID)
	if err != nil {
		logrus.Errorf("delete consul data for %s failed: %+v", nodeID, err)
		return err
	}

	// for server-added nodes
	err = svc.removeNodeFromServer(ctx, nodeID)
	if err != nil {
		logrus.Errorf("delete server data for %s failed: %+v", nodeID, err)
		return err
	}

	// remove prometheus data
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("delete metrics data for %s failed: %+v", nodeID, r)
			}
		}()

		err := svc.removeNodeFromPrometheus(ctx, nodeID)
		if err != nil {
			logrus.Errorf("delete metrics data for %s failed: %+v", nodeID, err)
		}
	}()

	return nil
}

func (svc *Service) removeNodeFromConsul(ctx context.Context, nodeID string) error {
	node, err := svc.consul.GetNode(nodeID)
	if err != nil {
		return err
	}
	if node == nil { // already gone
		return nil
	}

	_, err = svc.consul.DeregisterNode(node.Node.Node)
	return err
}

// RemoveService removes service of node,
// use id as service identifier if it id > 0,
// otherwise use name name as service identifier
func (svc *Service) RemoveService(ctx context.Context, nodeID string, id uint32, service string) error {
	err := svc.removeServiceFromQan(ctx, nodeID, service)
	if err != nil {
		logrus.Errorf("delete qan data for %s-%s failed: %+v", nodeID, service, err)
		return err
	}

	if id > 0 { // it's a server-added node
		err = svc.removeServiceFromServer(ctx, nodeID, service)
		if err != nil {
			logrus.Errorf("delete server data for %s-%s failed: %+v", nodeID, service, err)
			return err
		}
	} else { // it's a client-added node
		err = svc.removeServiceFromConsul(ctx, nodeID, service)
		if err != nil {
			logrus.Errorf("delete consul data for %s-%s failed: %+v", nodeID, service, err)
			return err
		}
	}

	// remove prometheus data
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("delete metrics data for %s-%s failed: %+v", nodeID, service, r)
			}
		}()

		err := svc.removeServiceFromPrometheus(ctx, nodeID, service)
		if err != nil {
			logrus.Errorf("delete metrics data for %s-%s failed: %+v", nodeID, service, err)
		}
	}()

	return nil
}

func (svc *Service) removeServiceFromPrometheus(ctx context.Context, nodeName, service string) error {
	deleteSeries := func(ctx context.Context, nodeName, service string) error {
		activeTargets, err := svc.prometheus.GetNodeServices(ctx)
		if err != nil {
			return fmt.Errorf("get active targets from prometheus for %s-%s failed: %s", nodeName, service, err.Error())
		}

		for _, target := range activeTargets {
			if target.Name == nodeName && string(target.Type) == service { // service has been re-added
				return nil
			}
		}

		queries := svc.genPromtheusQueries(nodeName, service)
		if queries != nil {
			err := svc.prometheus.DeleteSeries(queries)
			if err != nil {
				return fmt.Errorf("delete metrics data for %s failed: %s", nodeName, err.Error())
			}
		}

		return nil
	}

	deleteSeries(ctx, nodeName, service)

	// continually remove prometheus data incase there are some ongoing metrics task
	retryTimes := 35
	for i := 0; i < retryTimes; i++ {
		if i < 30 {
			<-time.Tick(1 * time.Second) // for second-level jobs
		} else {
			<-time.Tick(1 * time.Minute) // for minute-level jobs
		}

		err := deleteSeries(ctx, nodeName, service)
		if err != nil {
			logrus.Errorf("%s, try %d", err.Error(), i+1)
			continue
		}

		err = svc.prometheus.CleanTombstones(context.Background())
		if err != nil {
			logrus.Errorf("clean tombstones for %s failed: %s, try %d", nodeName, err.Error(), i+1)
			continue
		}

		return nil
	}

	return fmt.Errorf("delete metrics data for %s failed, tried %d times", nodeName, retryTimes)
}

func (svc *Service) removeServiceFromConsul(ctx context.Context, nodeID, name string) error {
	node, err := svc.consul.GetNode(nodeID)
	if err != nil {
		return err
	}
	if node == nil { // already gone
		return nil
	}

	deleteCount := 0
	for _, service := range node.Services {
		if service.Service != name {
			continue
		}

		_, err = svc.consul.DeregisterService(node.Node.Node, service.ID)
		if err != nil {
			return err
		}

		deleteCount++
	}

	if deleteCount >= len(node.Services) {
		_, err = svc.consul.DeregisterNode(node.Node.Node)
	}

	return err
}

func (svc *Service) removeServiceFromQan(ctx context.Context, nodeID, service string) error {
	var subsystemID int

	// server-side qan agent
	if service == string(models.MySQLdExporterAgentType) || service == string(models.MongoDBExporterAgentType) {
		return nil
	}

	if service == string(models.ClientMySQLQanAgentAgentType) {
		subsystemID = qan.SubsystemMySQL
	} else if service == string(models.ClientMongoDBQanAgentAgentType) {
		subsystemID = qan.SubsystemMongo
	}

	if subsystemID == 0 {
		// unknown qan service
		return nil
	}

	qanNodes, err := svc.GetQanNodes(ctx, nodeID, true)
	if err != nil {
		return err
	}

	if qanNodes == nil || len(qanNodes) == 0 {
		return nil
	}

	for _, node := range qanNodes {
		if node.Name != nodeID {
			continue
		}

		// not target
		if node.OSName == string(models.PMMServerNodeType) || node.SubsystemID != subsystemID {
			continue
		}

		// remove mysql qan queries
		agentUUID, err := svc.qan.GetAgentUUIDFromDB(ctx, nodeID, node.SubsystemID)
		if err != nil {
			return err
		}
		err = svc.qan.RemoveClientQAN(ctx, agentUUID, node.InstanceUUID)
		if err != nil {
			return err
		}

		go svc.removeQANData(ctx, nodeID, agentUUID, node.InstanceUUID)
	}

	return nil
}

func (svc *Service) removeServiceFromServer(ctx context.Context, nodeName, service string) error {
	removeNodeAndService := func(tx *reform.TX, agentNode *models.AgentNode, agentService *models.AgentService) error {
		if agentService == nil && agentNode == nil {
			return nil
		}
		var nodeID int32

		if agentNode != nil {
			nodeID = agentNode.NodeID
		}

		if agentService != nil {
			count, err := tx.Count(models.AgentServiceView, "WHERE service_id = ?", agentService.ServiceID)
			if err != nil {
				return errors.WithStack(err)
			}

			if count > 0 {
				return nil
			}

			if nodeID == 0 {
				var service models.Service
				err = tx.SelectOneTo(&service, "WHERE id = ?", agentService.ServiceID)
				if err != nil && err != sql.ErrNoRows {
					return errors.WithStack(err)
				}
				nodeID = service.NodeID
			}

			_, err = tx.DeleteFrom(models.ServiceTable, "WHERE id = ?", agentService.ServiceID)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		if nodeID == 0 {
			return nil
		}

		countAgent, err := tx.Count(models.AgentNodeView, "WHERE node_id = ?", nodeID)
		if err != nil {
			return errors.WithStack(err)
		}

		countService, err := tx.Count(models.ServiceTable, "WHERE node_id = ?", nodeID)
		if err != nil {
			return errors.WithStack(err)
		}

		if countAgent > 0 || countService > 0 {
			return nil
		}

		_, err = tx.DeleteFrom(models.NodeTable, "WHERE id = ?", nodeID)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	}

	return svc.db.InTransaction(func(tx *reform.TX) error {
		agentNode, err := models.AgentNodeByName(tx.Querier, nodeName, service)
		if err != nil {
			return errors.WithStack(err)
		}
		if agentNode != nil {
			_, err = tx.DeleteFrom(models.AgentNodeView, "WHERE node_id = ? AND agent_id = ?", agentNode.NodeID, agentNode.AgentID)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		agentService, err := models.AgentServiceByName(tx.Querier, nodeName, service)
		if err != nil {
			return errors.WithStack(err)
		}
		if agentService != nil {
			_, err = tx.DeleteFrom(models.AgentServiceView, "WHERE service_id = ? AND agent_id = ?", agentService.ServiceID, agentService.AgentID)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		if agentService == nil && agentNode == nil {
			return nil
		}

		// stop agents
		var agent models.Agent
		if agentNode != nil {
			err = tx.SelectOneTo(&agent, "WHERE id = ?", agentNode.AgentID)
		} else {
			err = tx.SelectOneTo(&agent, "WHERE id = ?", agentService.AgentID)
		}
		if err == sql.ErrNoRows {
			return nil
		}
		if err != nil {
			return errors.WithStack(err)
		}

		switch agent.Type {
		case models.MySQLdExporterAgentType:
			a := models.MySQLdExporter{ID: agent.ID}
			if err = tx.Reload(&a); err != nil {
				return errors.WithStack(err)
			}
			if svc.mysql.MySQLdExporterPath != "" {
				if err = svc.mysql.Supervisor.Stop(ctx, models.NameForSupervisor(a.Type, *a.ListenPort)); err != nil {
					return err
				}
			}

		case models.PostgresExporterAgentType:
			a := models.PostgresExporter{ID: agent.ID}
			if err = tx.Reload(&a); err != nil {
				return errors.WithStack(err)
			}
			if svc.postgresql.PostgresExporterPath != "" {
				if err = svc.postgresql.Supervisor.Stop(ctx, models.NameForSupervisor(a.Type, *a.ListenPort)); err != nil {
					return err
				}
			}

		case models.RDSExporterAgentType:
			a := models.RDSExporter{ID: agent.ID}
			if err = tx.Reload(&a); err != nil {
				return errors.WithStack(err)
			}
			if svc.rds.RDSExporterPath != "" {
				// update rds_exporter configuration
				config, err := svc.rds.UpdateRDSExporterConfig(tx)
				if err != nil {
					return err
				}

				// stop or restart rds_exporter
				name := models.NameForSupervisor(a.Type, *a.ListenPort)
				if err = svc.rds.Supervisor.Stop(ctx, name); err != nil {
					return err
				}
				if len(config.Instances) > 0 {
					if err = svc.rds.Supervisor.Start(ctx, svc.rds.RDSExporterServiceConfig(&a)); err != nil {
						return err
					}
				}
			}

		case models.QanAgentAgentType:
			a := models.QanAgent{ID: agent.ID}
			if err = tx.Reload(&a); err != nil {
				return errors.WithStack(err)
			}
			if svc.qan != nil {
				if err = svc.qan.RemoveMySQL(ctx, &a); err != nil {
					return err
				}
			}
		}

		// remove agent
		err = tx.Delete(&agent)
		if err != nil {
			return errors.WithStack(err)
		}

		err = removeNodeAndService(tx, agentNode, agentService)
		if err != nil {
			return errors.WithStack(err)
		}

		switch agent.Type {
		case models.RDSExporterAgentType:
			return svc.rds.ApplyPrometheusConfiguration(ctx, tx.Querier)
		case models.MySQLdExporterAgentType:
			return svc.mysql.ApplyPrometheusConfiguration(ctx, tx.Querier)
		case models.PostgresExporterAgentType:
			return svc.postgresql.ApplyPrometheusConfiguration(ctx, tx.Querier)
		}

		return nil
	})
}

func (svc *Service) removeNodeFromPrometheus(ctx context.Context, nodeID string) error {
	deleteSeries := func(ctx context.Context, nodeName string) error {
		activeTargets, err := svc.prometheus.GetNodeServices(ctx)
		if err != nil {
			return fmt.Errorf("get active targets from prometheus for %s failed: %s", nodeName, err.Error())
		}

		for _, target := range activeTargets {
			if string(target.Name) == nodeName { // service has been re-added
				return nil
			}
		}

		err = svc.prometheus.DeleteSeries(map[string]string{
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
			<-time.Tick(1 * time.Second) // for second-level jobs
		} else {
			<-time.Tick(1 * time.Minute) // for minute-level jobs
		}

		err := deleteSeries(ctx, nodeID)
		if err != nil {
			logrus.Errorf("delete metrics data for %s failed: %s, try %d", nodeID, err.Error(), i+1)
			continue
		}

		err = svc.prometheus.CleanTombstones(context.Background())
		if err != nil {
			logrus.Errorf("clean tombstones for %s failed: %s, try %d", nodeID, err.Error(), i+1)
			continue
		}

		return nil
	}

	return fmt.Errorf("delete metrics data for %s failed, tried %d times", nodeID, retryTimes)
}

func (svc *Service) removeNodeFromQan(ctx context.Context, nodeID string) error {
	qanNodes, err := svc.GetQanNodes(ctx, nodeID, true)
	if err != nil {
		return err
	}

	if qanNodes == nil || len(qanNodes) == 0 {
		return nil
	}

	for _, node := range qanNodes {
		if node.Name != nodeID {
			continue
		}

		if node.OSName == string(models.PMMServerNodeType) {
			continue
		}

		// remove mysql qan queries
		agentUUID, err := svc.qan.GetAgentUUIDFromDB(ctx, nodeID, node.SubsystemID)
		if err != nil {
			return err
		}

		err = svc.qan.RemoveClientQAN(ctx, agentUUID, node.InstanceUUID)
		if err != nil {
			return err
		}

		go svc.removeQANData(ctx, nodeID, agentUUID, node.InstanceUUID)
	}

	return nil
}

func (svc *Service) removeNodeFromServer(ctx context.Context, nodeID string) error {
	return svc.db.InTransaction(func(tx *reform.TX) error {
		dbNodes, err := tx.SelectAllFrom(models.NodeTable, "WHERE name = ?", nodeID)
		if err != nil {
			return errors.WithStack(err)
		}
		if len(dbNodes) == 0 {
			return nil
		}

		nodeIDs := make([]interface{}, len(dbNodes))
		for i, str := range dbNodes {
			nodeIDs[i] = str.(*models.Node).ID
		}

		dbServices, err := tx.FindAllFrom(models.ServiceTable, "node_id", nodeIDs...)
		if err != nil {
			return errors.WithStack(err)
		}
		serviceIDs := make([]interface{}, len(dbServices))
		for i, str := range dbServices {
			serviceIDs[i] = str.(*models.Service).ID
		}

		// remove associations of the service and agents
		var agentsForService []models.Agent
		if len(serviceIDs) > 0 {
			agentsForService, err = models.AgentsForServiceID(tx.Querier, serviceIDs...)
			if err != nil {
				return err
			}
			for i := range agentsForService {
				_, err = tx.DeleteFrom(
					models.AgentServiceView,
					"WHERE agent_id = ?",
					agentsForService[i].ID,
				)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}

		// remove associations of the node and agents
		agentsForNode, err := models.AgentsForNodeID(tx.Querier, nodeIDs...)
		if err != nil {
			return err
		}
		for i := range agentsForNode {
			_, err = tx.DeleteFrom(
				models.AgentNodeView,
				"WHERE agent_id = ?",
				agentsForNode[i].ID,
			)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		// stop agents
		agents := make(map[int32]models.Agent)
		for _, agent := range agentsForService {
			agents[agent.ID] = agent
		}
		for _, agent := range agentsForNode {
			agents[agent.ID] = agent
		}

		agentIDs := make([]interface{}, 0)
		for id, agent := range agents {
			agentIDs = append(agentIDs, id)

			switch agent.Type {
			case models.MySQLdExporterAgentType:
				a := models.MySQLdExporter{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.mysql.MySQLdExporterPath != "" {
					if err = svc.mysql.Supervisor.Stop(ctx, models.NameForSupervisor(a.Type, *a.ListenPort)); err != nil && !strings.Contains(err.Error(), services.ErrNoSuchFileOrDir.Error()) {
						return err
					}
				}

			case models.PostgresExporterAgentType:
				a := models.PostgresExporter{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.postgresql.PostgresExporterPath != "" {
					if err = svc.postgresql.Supervisor.Stop(ctx, models.NameForSupervisor(a.Type, *a.ListenPort)); err != nil && !strings.Contains(err.Error(), services.ErrNoSuchFileOrDir.Error()) {
						return err
					}
				}

			case models.RDSExporterAgentType:
				a := models.RDSExporter{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.rds.RDSExporterPath != "" {
					// update rds_exporter configuration
					config, err := svc.rds.UpdateRDSExporterConfig(tx)
					if err != nil {
						return err
					}

					// stop or restart rds_exporter
					name := models.NameForSupervisor(a.Type, *a.ListenPort)
					if err = svc.rds.Supervisor.Stop(ctx, name); err != nil && !strings.Contains(err.Error(), services.ErrNoSuchFileOrDir.Error()) {
						return err
					}
					if len(config.Instances) > 0 {
						if err = svc.rds.Supervisor.Start(ctx, svc.rds.RDSExporterServiceConfig(&a)); err != nil {
							return err
						}
					}
				}

			case models.QanAgentAgentType:
				a := models.QanAgent{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.qan != nil {
					agentUUID, err := svc.qan.GetAgentUUID()
					if err != nil {
						return err
					}

					<-time.Tick(1 * time.Second) // delay a little bit to avoid duplicate record in qan database
					if err = svc.qan.RemoveMySQL(ctx, &a); err != nil {
						return err
					}

					go svc.removeQANData(ctx, nodeID, agentUUID, *a.QANDBInstanceUUID)
				}
			}
		}

		// remove agents
		for _, agentID := range agentIDs {
			_, err = tx.DeleteFrom(models.AgentTable, "WHERE id = ?", agentID)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		// delete services
		for _, serviceID := range serviceIDs {
			_, err = tx.DeleteFrom(models.ServiceTable, "WHERE id = ?", serviceID)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		// delete nodes
		for _, nodeID := range nodeIDs {
			_, err = tx.DeleteFrom(models.NodeTable, "WHERE id = ?", nodeID)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		// reconfigure mysql prometheus
		err = svc.mysql.ApplyPrometheusConfiguration(ctx, tx.Querier)
		if err != nil {
			return errors.WithStack(err)
		}

		// reconfigure postgresql prometheus
		err = svc.postgresql.ApplyPrometheusConfiguration(ctx, tx.Querier)
		if err != nil {
			return errors.WithStack(err)
		}

		// // reconfigure rds prometheus
		err = svc.rds.ApplyPrometheusConfiguration(ctx, tx.Querier)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}

// GetConsulNodes returns client nodes from consul
func (svc *Service) GetConsulNodes(ctx context.Context) ([]ClientNode, error) {
	var clientNodes []ClientNode

	nodes, err := svc.consul.GetNodes()
	if err != nil {
		logger.Get(ctx).Errorf("get nodes from consul failed: %+v", err)
		return nil, err
	}
	for _, node := range nodes {
		cNode, err := svc.consul.GetNode(node.Node)
		if err != nil {
			logger.Get(ctx).Errorf("get consul services from node failed: %+v", err)
			continue
		}

		isReserved := false
		services := make([]ClientNodeService, 0)
		for _, service := range cNode.Services {
			if service.ID == "consul" && service.Service == "consul" { // reserved service
				isReserved = true
				break
			}

			cns := ClientNodeService{
				Name:    service.Service,
				Address: node.Address,
				Port:    service.Port,
				Region:  string(models.ClientNodeRegion),
			}
			for _, tag := range service.Tags {
				if strings.HasPrefix(tag, "distro_") {
					parts := strings.SplitN(tag, "_", 2)
					if len(parts) > 1 {
						cns.Distro = parts[1]
					}
				} else if strings.HasPrefix(tag, "version_") {
					parts := strings.SplitN(tag, "_", 2)
					if len(parts) > 1 {
						cns.Version = parts[1]
					}
				}
			}
			services = append(services, cns)
		}

		if isReserved {
			continue
		}

		clientNodes = append(clientNodes, ClientNode{
			Name:     node.Node,
			Services: services,
		})
	}

	return clientNodes, nil
}

// GetQanNodes returns client nodes from QAN
func (svc *Service) GetQanNodes(ctx context.Context, name string, checkData bool) ([]qan.UnremovedNode, error) {
	return svc.qan.GetUnremovedNodes(ctx, name, checkData)
}

// GetPrometheusNodes returns client nodes from prometheus
func (svc *Service) GetPrometheusNodes(ctx context.Context) ([]prometheus.NodeService, error) {
	return svc.prometheus.GetNodeServices(ctx)
}

// GetRegionFromAgentType returns region of agent type
func (svc *Service) GetRegionFromAgentType(agentType models.AgentType) string {
	switch agentType {
	case models.MySQLdExporterAgentType, models.PostgresExporterAgentType,
		models.RDSExporterAgentType, models.QanAgentAgentType,
		models.NodeExporterAgentType, models.ProxySQLExporterAgentType,
		models.MongoDBExporterAgentType:
		return string(models.RemoteNodeRegion)
	case models.ClientNodeExporterAgentType, models.ClientMySQLdExporterAgentType,
		models.ClientMySQLQanAgentAgentType, models.ClientMongoDBExporterAgentType,
		models.ClientMongoDBQanAgentAgentType, models.ClientPostgresExporterAgentType,
		models.ClientProxySQLExporterAgentType:
		return string(models.ClientNodeRegion)
	}

	return ""
}

func (svc *Service) genPromtheusQueries(nodeName string, service string) map[string]string {
	queries := map[string]string{
		"instance=": nodeName,
	}

	agentType := models.AgentType(service)
	switch agentType {
	case models.MySQLdExporterAgentType, models.PostgresExporterAgentType,
		models.NodeExporterAgentType, models.ProxySQLExporterAgentType,
		models.MongoDBExporterAgentType:
		queries["region="] = string(models.RemoteNodeRegion)
	case models.ClientNodeExporterAgentType, models.ClientMySQLdExporterAgentType,
		models.ClientMongoDBExporterAgentType, models.ClientPostgresExporterAgentType,
		models.ClientProxySQLExporterAgentType:
		queries["region="] = string(models.ClientNodeRegion)
	case models.RDSExporterAgentType:
		break
	default:
		return nil
	}

	switch agentType {
	case models.MySQLdExporterAgentType, models.ClientMySQLdExporterAgentType:
		queries["job="] = "mysql"
	case models.PostgresExporterAgentType, models.ClientPostgresExporterAgentType:
		queries["job="] = "postgresql"
	case models.MongoDBExporterAgentType, models.ClientMongoDBExporterAgentType:
		queries["job="] = "mongodb"
	case models.NodeExporterAgentType, models.ClientNodeExporterAgentType:
		queries["job="] = "linux"
	case models.ProxySQLExporterAgentType, models.ClientProxySQLExporterAgentType:
		queries["job="] = "proxysql"
	case models.RDSExporterAgentType:
		queries["job=~"] = "rds-*"
	default:
		return nil
	}

	return queries
}

func (svc *Service) removeQANData(ctx context.Context, nodeID, agentUUID, instanceUUID string) {
	deleteData := func() error {
		nodes, err := svc.GetQanNodes(ctx, nodeID, false)
		if err != nil {
			return err
		}

		for _, node := range nodes {
			if node.InstanceUUID == instanceUUID { // qan node re-added
				return nil
			}
		}

		return svc.qan.RemoveQANData(ctx, agentUUID, instanceUUID)
	}

	deleteData()
	// continually remove qan data incase there are some ongoing qan task
	retryTimes := 30
	for i := 0; i < retryTimes; i++ {
		<-time.Tick(1 * time.Second)

		err := deleteData()
		if err != nil {
			logrus.Errorf("remove qan data for %s-%s-%s failed: %+v", nodeID, agentUUID, instanceUUID, err)
		}
	}
}
