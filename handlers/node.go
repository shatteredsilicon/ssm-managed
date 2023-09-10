package handlers

import (
	"context"
	"strconv"
	"strings"

	"github.com/shatteredsilicon/ssm-managed/api"
	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/services/node"
	nodeSvc "github.com/shatteredsilicon/ssm-managed/services/node"
	"github.com/shatteredsilicon/ssm-managed/services/prometheus"
	"github.com/shatteredsilicon/ssm-managed/services/qan"
	"github.com/shatteredsilicon/ssm-managed/services/remote"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
)

// NodeServer server for node apis
type NodeServer struct {
	Node   *nodeSvc.Service
	Remote *remote.Service
}

// Remove removes nodes.
func (s *NodeServer) Remove(ctx context.Context, req *api.NodeRemoveRequest) (*api.NodeRemoveResponse, error) {
	if err := s.Node.RemoveNode(ctx, req.Name); err != nil {
		logger.Get(ctx).Errorf("remove node %s failed: %+v", req.Name, err)
		return nil, err
	}

	var resp api.NodeRemoveResponse
	return &resp, nil
}

// RemoveService removes specify service of node
func (s *NodeServer) RemoveService(ctx context.Context, req *api.NodeRemoveServiceRequest) (*api.NodeRemoveServiceResponse, error) {
	if err := s.Node.RemoveService(ctx, req.Name, req.Id, req.Type); err != nil {
		logger.Get(ctx).Errorf("remove service %s-%s failed: %+v", req.Name, req.Type, err)
		return nil, err
	}

	var resp api.NodeRemoveServiceResponse
	return &resp, nil
}

// List returns client-added and server-added nodes
func (s *NodeServer) List(ctx context.Context, req *api.NodeListRequest) (*api.NodeListResponse, error) {
	resp := api.NodeListResponse{
		Instances: []*api.NodeInstance{},
	}

	res, err := s.Remote.ListFull(ctx)
	if err != nil {
		logger.Get(ctx).Errorf("%+v", err)
		return nil, err
	}
	resp.Instances = s.putRemoteNodes(res, resp.Instances)

	nodes, err := s.Node.GetConsulNodes(ctx)
	if err != nil {
		logger.Get(ctx).Errorf("get client nodes failed: %+v", err)
		return nil, err
	}
	resp.Instances = s.putConsulNodes(nodes, resp.Instances)

	qanNodes, err := s.Node.GetQanNodes(ctx, "", true)
	if err != nil {
		logger.Get(ctx).Errorf("get unremoved qan nodes failed: %+v", err)
		return nil, err
	}
	resp.Instances = s.putQanNodes(qanNodes, resp.Instances)

	promNodes, err := s.Node.GetPrometheusNodes(ctx)
	if err != nil {
		logger.Get(ctx).Errorf("get unremoved prometheus nodes failed: %+v", err)
		return nil, err
	}
	resp.Instances = s.putPrometheusNodes(promNodes, resp.Instances)

	return &resp, nil
}

func (s *NodeServer) putRemoteNodes(nodes []remote.FullInstance, respNodes []*api.NodeInstance) []*api.NodeInstance {
	for _, db := range nodes {
		services := make([]*api.NodeInstanceService, len(db.Service.Agents))
		for i, agent := range db.Service.Agents {
			services[i] = &api.NodeInstanceService{
				Id:            uint32(db.Service.ID),
				Type:          string(agent.Type),
				Address:       *db.Service.Address,
				Port:          uint32(*db.Service.Port),
				Region:        db.Node.Region,
				Engine:        *db.Service.Engine,
				EngineVersion: *db.Service.EngineVersion,
			}
		}

		found := false
		for i := range respNodes {
			if respNodes[i].Name == db.Node.Name {
				respNodes[i].Services = append(respNodes[i].Services, services...)
				found = true
				break
			}
		}
		if found {
			continue
		}

		respNodes = append(respNodes, &api.NodeInstance{
			Name:     db.Node.Name,
			Services: services,
		})
	}
	return respNodes
}

func (s *NodeServer) putConsulNodes(nodes []node.ClientNode, respNodes []*api.NodeInstance) []*api.NodeInstance {
	for _, node := range nodes {
		services := make([]*api.NodeInstanceService, len(node.Services))
		for i, service := range node.Services {
			services[i] = &api.NodeInstanceService{
				Type:          service.Name,
				Address:       service.Address,
				Port:          uint32(service.Port),
				Region:        service.Region,
				Engine:        service.Distro,
				EngineVersion: service.Version,
			}
		}

		found := false
		for i := range respNodes {
			if respNodes[i].Name == node.Name {
				respNodes[i].Services = append(respNodes[i].Services, services...)
				found = true
				break
			}
		}
		if found {
			continue
		}

		respNodes = append(respNodes, &api.NodeInstance{
			Name:     node.Name,
			Services: services,
		})
	}
	return respNodes
}

func (s *NodeServer) putQanNodes(nodes []qan.UnremovedNode, respNodes []*api.NodeInstance) []*api.NodeInstance {
	for _, node := range nodes {
		if node.Name == string(models.SSMServerNodeType) { // server itself
			continue
		}

		var nis api.NodeInstanceService
		if node.OSName == string(models.SSMServerNodeType) {
			nis.Region = string(models.RemoteNodeRegion)
			nis.Type = string(models.QanAgentAgentType)
		} else {
			nis.Region = string(models.ClientNodeRegion)
			nis.Type = string(models.ClientMySQLQanAgentAgentType)
			if node.SubsystemID == qan.SubsystemMongo {
				nis.Type = string(models.ClientMongoDBQanAgentAgentType)
			}
		}

		exists := false
		for i := range respNodes {
			if respNodes[i].Name != node.Name {
				continue
			}

			serviceExists := false
			for _, service := range respNodes[i].Services {
				if node.SubsystemID == qan.SubsystemMySQL && (service.Type == string(models.QanAgentAgentType) ||
					service.Type == string(models.ClientMySQLQanAgentAgentType)) {
					serviceExists = true
					break
				} else if node.SubsystemID == qan.SubsystemMongo && (service.Type == string(models.QanAgentAgentType) ||
					service.Type == string(models.ClientMongoDBQanAgentAgentType)) {
					serviceExists = true
					break
				}
			}

			if !serviceExists {
				respNodes[i].Services = append(respNodes[i].Services, &nis)
			}

			exists = true
			break
		}

		if exists { // already recorded
			continue
		}

		respNodes = append(respNodes, &api.NodeInstance{
			Name:     node.Name,
			Services: []*api.NodeInstanceService{&nis},
		})
	}

	return respNodes
}

func (s *NodeServer) putPrometheusNodes(nodes []prometheus.NodeService, respNodes []*api.NodeInstance) []*api.NodeInstance {
	for _, node := range nodes {
		if node.Name == string(models.SSMServerNodeType) { // server itself
			continue
		}

		addressParts := strings.Split(node.Endpoint, ":")
		nis := api.NodeInstanceService{
			Type:    string(node.Type),
			Region:  s.Node.GetRegionFromAgentType(node.Type),
			Address: addressParts[0],
		}
		if len(addressParts) > 1 {
			port, _ := strconv.Atoi(addressParts[1])
			nis.Port = uint32(port)
		}

		exists := false
		for i := range respNodes {
			if respNodes[i].Name != node.Name {
				continue
			}

			serviceExists := false
			for _, service := range respNodes[i].Services {
				if service.Type == nis.Type {
					serviceExists = true
					break
				}
			}

			if !serviceExists {
				respNodes[i].Services = append(respNodes[i].Services, &nis)
			}

			exists = true
			break
		}

		if exists { // already recorded
			continue
		}

		respNodes = append(respNodes, &api.NodeInstance{
			Name:     node.Name,
			Services: []*api.NodeInstanceService{&nis},
		})
	}

	return respNodes
}
