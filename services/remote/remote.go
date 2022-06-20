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
