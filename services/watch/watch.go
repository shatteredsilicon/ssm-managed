package watch

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/services/consul"
	"github.com/shatteredsilicon/ssm-managed/services/mysql"
	"github.com/shatteredsilicon/ssm-managed/services/node"
	"github.com/shatteredsilicon/ssm-managed/services/postgresql"
	"github.com/shatteredsilicon/ssm-managed/services/rds"
	"github.com/shatteredsilicon/ssm-managed/services/remote"
	"github.com/shatteredsilicon/ssm-managed/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

const (
	engineWatchInterval = time.Hour
)

// Service metric service
type Service struct {
	node       *node.Service
	remote     *remote.Service
	rds        *rds.Service
	mysql      *mysql.Service
	postgresql *postgresql.Service
	db         *reform.DB
	consul     *consul.Client
	logger     *logrus.Entry
}

// NewService returns a new metric service
func NewService(
	node *node.Service,
	remote *remote.Service,
	rds *rds.Service,
	mysql *mysql.Service,
	postgresql *postgresql.Service,
	db *reform.DB,
	consul *consul.Client,
	logger *logrus.Entry,
) *Service {
	return &Service{
		node:       node,
		remote:     remote,
		db:         db,
		rds:        rds,
		mysql:      mysql,
		postgresql: postgresql,
		consul:     consul,
		logger:     logger,
	}
}

// Run runs the metric service
func (s *Service) Run(ctx context.Context) {
	ticker := time.NewTicker(engineWatchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.updateEngine(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) updateEngine(ctx context.Context) {
	l := s.logger.WithField("function", "updateEngine")

	remoteInstances, err := s.remote.ListFull(ctx)
	if err != nil {
		l.Errorf("Failed to get list of remote nodes: %s\n", err)
		return
	}

	rdsInstances := make(map[string][]rds.Instance)
	for _, remoteInstance := range remoteInstances {
		var oldEngine, oldVersion, newEngine, newVersion string
		if remoteInstance.Service.Engine != nil {
			oldEngine = *remoteInstance.Service.Engine
		}
		if remoteInstance.Service.EngineVersion != nil {
			oldVersion = *remoteInstance.Service.EngineVersion
		}

		if remoteInstance.Service.Type == models.RDSServiceType {
			var instances []rds.Instance
			var accessKey, secretKey string
			if remoteInstance.Service.AWSAccessKey != nil {
				accessKey = *remoteInstance.Service.AWSAccessKey
			}
			if remoteInstance.Service.AWSSecretKey != nil {
				secretKey = *remoteInstance.Service.AWSSecretKey
			}

			key := accessKey + ":" + secretKey
			if _, ok := rdsInstances[key]; ok {
				instances = rdsInstances[key]
			} else {
				instances, err = s.rds.Discover(ctx, accessKey, secretKey)
				if err != nil {
					l.Errorf("Got an error when discovering rds service: %s\n", err)

					// ignore error if there are some results
					if len(instances) == 0 {
						continue
					}
				}

				rdsInstances[key] = instances
			}

			var rdsInstance *rds.Instance
			for _, inst := range instances {
				if remoteInstance.Node.Name == inst.Node.Name && remoteInstance.Node.Region == inst.Node.Region {
					rdsInstance = &inst
					break
				}
			}
			if rdsInstance == nil {
				continue
			}

			if rdsInstance.Service.Engine != nil {
				newEngine = *rdsInstance.Service.Engine
			}
			if rdsInstance.Service.EngineVersion != nil {
				newVersion = *rdsInstance.Service.EngineVersion
			}
		} else if remoteInstance.Service.Type == models.MySQLServiceType || remoteInstance.Service.Type == models.PostgreSQLServiceType {
			var address, username, password string
			var port uint32
			if remoteInstance.Service.Address != nil {
				address = *remoteInstance.Service.Address
			}
			if remoteInstance.Service.Port != nil {
				port = uint32(*remoteInstance.Service.Port)
			}
			// get service username/password from agents
			for _, agent := range remoteInstance.Service.Agents {
				if username == "" && agent.ServiceUsername != nil {
					username = *agent.ServiceUsername
				}
				if password == "" && agent.ServicePassword != nil {
					password = *agent.ServicePassword
				}
				if username != "" && password != "" {
					break
				}
			}
			if remoteInstance.Service.Type == models.MySQLServiceType {
				newEngine, newVersion, err = s.mysql.EngineAndEngineVersion(ctx, address, port, username, password)
				if err != nil {
					l.Errorf("Got an error when getting engine info for mysql service: %s\n", err)
					continue
				}
			} else if remoteInstance.Service.Type == models.PostgreSQLServiceType {
				newEngine, newVersion, err = s.postgresql.EngineAndEngineVersion(ctx, address, port, username, password)
				if err != nil {
					l.Errorf("Got an error when getting engine info for postgresql service: %s\n", err)
					continue
				}
			}
		} else {
			continue
		}

		// engine version not changed, skip
		if oldEngine == newEngine && oldVersion == newVersion {
			continue
		}

		_, err := s.db.ExecContext(
			ctx,
			"UPDATE "+models.RemoteServiceTable.Name()+" SET engine = ?, engine_version = ? WHERE id = ?",
			newEngine,
			newVersion,
			remoteInstance.Service.ID,
		)
		if err != nil {
			l.Errorf("Failed to update service engine info for instance '%s', err: %s\n", remoteInstance.Node.Name, err)
			continue
		}
	}

	clientInstances, err := s.node.GetConsulNodes(ctx)
	if err != nil {
		l.Errorf("Failed to get list of client nodes: %s\n", err)
		return
	}

	instanceServices := make(map[string][]string)
	for _, instance := range clientInstances {
		if len(instance.Services) > 0 {
			instanceServices[instance.Name] = make([]string, 0)
		}
		for _, service := range instance.Services {
			serviceType := strings.SplitN(service.Name, ":", 2)[0]
			if utils.SliceContains(instanceServices[instance.Name], serviceType) {
				continue
			}

			instanceServices[instance.Name] = append(instanceServices[instance.Name], serviceType)
		}
	}

	serviceEngines, err := s.node.GetServiceEngine(instanceServices)
	if err != nil {
		l.Errorf("Failed to get service engine info: %s\n", err)
		return
	}

	for _, instance := range clientInstances {
		ses, ok := serviceEngines[instance.Name]
		if !ok {
			continue
		}

		for _, service := range instance.Services {
			serviceType := strings.SplitN(service.Name, ":", 2)[0]
			se, ok := ses[serviceType]
			if !ok {
				continue
			}

			if se.Engine == service.Distro && se.Version == service.Version {
				continue
			}

			cNode, err := s.consul.GetNode(instance.Name)
			if err != nil {
				l.Errorf("Failed to get consul node for instance '%s', err: %s\n", instance.Name, err)
				continue
			}

			var consulService *api.AgentService
			for _, consulSvc := range cNode.Services {
				if consulSvc.Service != service.Name {
					continue
				}

				consulService = consulSvc
				for tagI, tag := range consulSvc.Tags {
					if strings.HasPrefix(tag, "distro_") && se.Engine != "" && se.Engine != service.Distro {
						consulSvc.Tags[tagI] = "distro_" + se.Engine
					}
					if strings.HasPrefix(tag, "version_") && se.Version != "" && se.Version != service.Version {
						consulSvc.Tags[tagI] = "version_" + se.Version
					}
				}
			}

			if consulService == nil {
				continue
			}

			_, err = s.consul.Register(&api.CatalogRegistration{
				ID:              cNode.Node.ID,
				Node:            cNode.Node.Node,
				Address:         cNode.Node.Address,
				TaggedAddresses: cNode.Node.TaggedAddresses,
				NodeMeta:        cNode.Node.Meta,
				Datacenter:      cNode.Node.Datacenter,
				Service:         consulService,
			}, nil)
			if err != nil {
				l.Errorf("Failed to update consul service for instance '%s', err: %s\n", instance.Name, err)
				continue
			}
		}
	}
}
