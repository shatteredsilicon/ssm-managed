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

//reform:agent_nodes
type AgentNode struct {
	AgentID int32 `reform:"agent_id"`
	NodeID  int32 `reform:"node_id"`
}

// AgentsForNodeID returns agents providing insights for a given node.
func AgentsForNodeID(q *reform.Querier, nodeIDs ...interface{}) ([]Agent, error) {
	agentNodes, err := q.FindAllFrom(AgentNodeView, "node_id", nodeIDs...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	agentIDs := make([]interface{}, len(agentNodes))
	for i, str := range agentNodes {
		agentIDs[i] = str.(*AgentNode).AgentID
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

// AgentNodeByName returns agent_nodes
func AgentNodeByName(q *reform.Querier, nodeName, agentType string) (*AgentNode, error) {
	var agentID, nodeID int32
	err := q.QueryRow(`
SELECT agn.agent_id, agn.node_id
FROM agent_nodes agn
JOIN nodes ON agn.node_id = nodes.id
JOIN agents ON agn.agent_id = agents.id
WHERE nodes.name = ? and agents.type = ?
`, nodeName, agentType).Scan(&agentID, &nodeID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &AgentNode{
		AgentID: agentID,
		NodeID:  nodeID,
	}, nil
}
