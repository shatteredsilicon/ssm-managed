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

package models

import (
	"database/sql"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

//go:generate reform

//reform:agent_services
type AgentService struct {
	AgentID   int32 `reform:"agent_id"`
	ServiceID int32 `reform:"service_id"`
}

// AgentServiceDetail AgentService with detail
type AgentServiceDetail struct {
	AgentService
	NodeID            int32  `reform:"node_id"`
	NodeType          string `reform:"node_type"`
	QanDBInstanceUUID string `reform:"qan_db_instance_uuid"`
}

// AgentsForServiceID returns agents providing insights for given services.
func AgentsForServiceID(q *reform.Querier, serviceIDs ...interface{}) ([]Agent, error) {
	agentServices, err := q.FindAllFrom(AgentServiceView, "service_id", serviceIDs...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	agentIDs := make([]interface{}, len(agentServices))
	for i, str := range agentServices {
		agentIDs[i] = str.(*AgentService).AgentID
	}

	if len(agentIDs) == 0 {
		return []Agent{}, nil
	}

	structs, err := q.FindAllFrom(AgentTable, "id", agentIDs...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	agents := make([]Agent, len(structs))
	for i, str := range structs {
		agents[i] = *str.(*Agent)
	}
	return agents, nil
}

// AgentServiceByName returns agent_service
func AgentServiceByName(q *reform.Querier, nodeName, agentType string) (*AgentServiceDetail, error) {
	var agentID, serviceID, nodeID int32
	var nodeType string
	var qanDBInstanceUUID string
	err := q.QueryRow(`
SELECT agsv.agent_id, agsv.service_id, nodes.type AS node_type, nodes.id AS node_id, agents.qan_db_instance_uuid
FROM agent_services agsv
JOIN services ON agsv.service_id = services.id
JOIN nodes ON services.node_id = nodes.id
JOIN agents ON agsv.agent_id = agents.id
WHERE nodes.name = ? and agents.type = ?
`, nodeName, agentType).Scan(&agentID, &serviceID, &nodeType, &nodeID, &qanDBInstanceUUID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &AgentServiceDetail{
		AgentService: AgentService{
			AgentID:   agentID,
			ServiceID: serviceID,
		},
		NodeID:            nodeID,
		NodeType:          nodeType,
		QanDBInstanceUUID: qanDBInstanceUUID,
	}, nil
}
