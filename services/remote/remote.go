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

// Package remote contains business logic of working with all Remote instances.
package remote

import (
	"context"

	"gopkg.in/reform.v1"

	"github.com/shatteredsilicon/ssm-managed/models"
)

// ServiceConfig contains configuration for remote.Service
type ServiceConfig struct {
	DB *reform.DB
}

// Service is responsible for interactions with Remote instances.
type Service struct {
	*ServiceConfig
}

// NewService creates a new service.
func NewService(config *ServiceConfig) (*Service, error) {
	return &Service{
		ServiceConfig: config,
	}, nil
}

// Instance contains data about node and service placed on the node
type Instance struct {
	Node    models.RemoteNode
	Service models.RemoteService
}

// FullInstanceService contains service and agents placed on the service
type FullInstanceService struct {
	models.FullService
	Agents []models.FullAgent
}

// FullInstance contains data about node, service and agent placed on the node
type FullInstance struct {
	Node    models.RemoteNode
	Service FullInstanceService
}

// List returns a list of all remote nodes (including RDS nodes).
func (svc *Service) List(ctx context.Context) ([]Instance, error) {
	var res []Instance
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		structs, e := tx.SelectAllFrom(models.RemoteNodeTable, "WHERE type IN (?, ?) ORDER BY id", models.RDSNodeType, models.RemoteNodeType)
		if e != nil {
			return e
		}
		nodes := make([]models.RemoteNode, len(structs))
		for i, str := range structs {
			nodes[i] = *str.(*models.RemoteNode)
		}

		structs, e = tx.SelectAllFrom(models.RemoteServiceTable, "ORDER BY id")
		if e != nil {
			return e
		}
		services := make([]models.RemoteService, len(structs))
		for i, str := range structs {
			services[i] = *str.(*models.RemoteService)
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

// ListFull returns a list of all remote nodes (including RDS nodes),
// with all services and agents data.
func (svc *Service) ListFull(ctx context.Context) ([]FullInstance, error) {
	var res []FullInstance
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		structs, e := tx.SelectAllFrom(models.RemoteNodeTable, "WHERE type IN (?, ?) ORDER BY id", models.RDSNodeType, models.RemoteNodeType)
		if e != nil {
			return e
		}
		nodes := make([]models.RemoteNode, len(structs))
		for i, str := range structs {
			nodes[i] = *str.(*models.RemoteNode)
		}

		structs, e = tx.SelectAllFrom(models.FullServiceTable, "ORDER BY id")
		if e != nil {
			return e
		}
		services := make([]models.FullService, len(structs))
		for i, str := range structs {
			services[i] = *str.(*models.FullService)
		}

		structs, e = tx.SelectAllFrom(models.FullAgentTable, "ORDER BY id")
		if e != nil {
			return e
		}
		agents := make(map[int32]models.FullAgent, len(structs))
		for _, str := range structs {
			agent := *str.(*models.FullAgent)
			agents[agent.ID] = agent
		}

		structs, e = tx.SelectAllFrom(models.AgentServiceView, "")
		if e != nil {
			return e
		}
		serviceAgents := make(map[int32][]models.FullAgent)
		for _, str := range structs {
			agentService := *str.(*models.AgentService)
			agent, ok := agents[agentService.AgentID]
			if !ok {
				continue
			}

			if _, ok := serviceAgents[agentService.ServiceID]; ok {
				serviceAgents[agentService.ServiceID] = append(serviceAgents[agentService.ServiceID], agent)
			} else {
				serviceAgents[agentService.ServiceID] = []models.FullAgent{agent}
			}
		}

		for _, node := range nodes {
			for _, service := range services {
				if node.ID == service.NodeID {
					resService := FullInstance{
						Node: node,
						Service: FullInstanceService{
							FullService: service,
						},
					}

					if agents, ok := serviceAgents[service.ID]; ok {
						resService.Service.Agents = agents
					}

					res = append(res, resService)
				}
			}
		}
		return nil
	})
	return res, err
}
