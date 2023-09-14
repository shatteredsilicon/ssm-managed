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
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/reform.v1"
)

//go:generate reform

const (
	// maximum time for connecting to the database
	sqlDialTimeout = 5 * time.Second
)

type AgentType string

// AgentType agent types for exporters and agents
const (
	MySQLdExporterAgentType         AgentType = "mysqld_exporter"
	PostgresExporterAgentType       AgentType = "postgres_exporter"
	RDSExporterAgentType            AgentType = "rds_exporter"
	QanAgentAgentType               AgentType = "qan-agent"
	NodeExporterAgentType           AgentType = "node_exporter"
	ProxySQLExporterAgentType       AgentType = "proxysql_exporter"
	MongoDBExporterAgentType        AgentType = "mongodb_exporter"
	SNMPExporterAgentType           AgentType = "snmp_exporter"
	ClientNodeExporterAgentType     AgentType = "linux:metrics"
	ClientMySQLdExporterAgentType   AgentType = "mysql:metrics"
	ClientMySQLQanAgentAgentType    AgentType = "mysql:queries"
	ClientMongoDBExporterAgentType  AgentType = "mongodb:metrics"
	ClientMongoDBQanAgentAgentType  AgentType = "mongodb:queries"
	ClientPostgresExporterAgentType AgentType = "postgresql:metrics"
	ClientProxySQLExporterAgentType AgentType = "proxysql:metrics"
)

// NameForSupervisor returns a name of agent for supervisor.
func NameForSupervisor(typ AgentType, listenPort uint16) string {
	return fmt.Sprintf("ssm-%s-%d", typ, listenPort)
}

//reform:agents
type Agent struct {
	ID                int32     `reform:"id,pk"`
	Type              AgentType `reform:"type"`
	RunsOnNodeID      int32     `reform:"runs_on_node_id"`
	QanDBInstanceUUID *string   `reform:"qan_db_instance_uuid"`

	// TODO Does it really belong there? Remove when we have agent without one.
	ListenPort *uint16 `reform:"listen_port"`
}

//reform:agents
type MySQLdExporter struct {
	ID           int32     `reform:"id,pk"`
	Type         AgentType `reform:"type"`
	RunsOnNodeID int32     `reform:"runs_on_node_id"`

	ServiceUsername        *string `reform:"service_username"`
	ServicePassword        *string `reform:"service_password"`
	ListenPort             *uint16 `reform:"listen_port"`
	MySQLDisableTablestats *bool   `reform:"mysql_disable_tablestats"`
}

func (m *MySQLdExporter) DSN(service *MySQLService) string {
	cfg := mysql.NewConfig()
	cfg.User = *m.ServiceUsername
	cfg.Passwd = *m.ServicePassword

	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(*service.Address, strconv.Itoa(int(*service.Port)))

	cfg.Timeout = sqlDialTimeout

	// TODO TLSConfig: "true", https://jira.percona.com/browse/PMM-1727
	// TODO Other parameters?
	return cfg.FormatDSN()
}

// binary name is postgres_exporter, that's why PostgresExporter below is not PostgreSQLExporter

// PostgresExporter exports PostgreSQL metrics.
//
//reform:agents
type PostgresExporter struct {
	ID           int32     `reform:"id,pk"`
	Type         AgentType `reform:"type"`
	RunsOnNodeID int32     `reform:"runs_on_node_id"`

	ServiceUsername *string `reform:"service_username"`
	ServicePassword *string `reform:"service_password"`
	ListenPort      *uint16 `reform:"listen_port"`
}

// DSN returns DSN for PostgreSQL service.
func (p *PostgresExporter) DSN(service *PostgreSQLService) string {
	q := make(url.Values)
	q.Set("sslmode", "disable") // TODO https://jira.percona.com/browse/PMM-1727
	q.Set("connect_timeout", strconv.Itoa(int(sqlDialTimeout.Seconds())))

	address := net.JoinHostPort(*service.Address, strconv.Itoa(int(*service.Port)))
	uri := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(*p.ServiceUsername, *p.ServicePassword),
		Host:     address,
		Path:     "postgres",
		RawQuery: q.Encode(),
	}
	return uri.String()
}

//reform:agents
type RDSExporter struct {
	ID           int32     `reform:"id,pk"`
	Type         AgentType `reform:"type"`
	RunsOnNodeID int32     `reform:"runs_on_node_id"`

	ListenPort *uint16 `reform:"listen_port"`
}

//reform:agents
type QanAgent struct {
	ID           int32     `reform:"id,pk"`
	Type         AgentType `reform:"type"`
	RunsOnNodeID int32     `reform:"runs_on_node_id"`

	ServiceUsername   *string `reform:"service_username"`
	ServicePassword   *string `reform:"service_password"`
	ListenPort        *uint16 `reform:"listen_port"`
	QANDBInstanceUUID *string `reform:"qan_db_instance_uuid"` // MySQL instance UUID in QAN
}

//reform:agents
type NodeExporter struct {
	ID           int32     `reform:"id,pk"`
	Type         AgentType `reform:"type"`
	RunsOnNodeID int32     `reform:"runs_on_node_id"`

	ListenPort *uint16 `reform:"listen_port"`
}

// reform:agents
type FullAgent struct {
	ID                int32     `reform:"id,pk"`
	Type              AgentType `reform:"type"`
	RunsOnNodeID      int32     `reform:"runs_on_node_id"`
	QanDBInstanceUUID *string   `reform:"qan_db_instance_uuid"`

	ServiceUsername        *string `reform:"service_username"`
	ServicePassword        *string `reform:"service_password"`
	ListenPort             *uint16 `reform:"listen_port"`
	MySQLDisableTablestats *bool   `reform:"mysql_disable_tablestats"`
}

func (q *QanAgent) DSN(service *MySQLService) string {
	cfg := mysql.NewConfig()
	cfg.User = *q.ServiceUsername
	cfg.Passwd = *q.ServicePassword

	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(*service.Address, strconv.Itoa(int(*service.Port)))

	cfg.Timeout = sqlDialTimeout

	// TODO TLSConfig: "true", https://jira.percona.com/browse/PMM-1727
	// TODO Other parameters?
	return cfg.FormatDSN()
}

// SNMPExporter exports SNMP metrics.
//
//reform:agents
type SNMPExporter struct {
	ID           int32     `reform:"id,pk"`
	Type         AgentType `reform:"type"`
	RunsOnNodeID int32     `reform:"runs_on_node_id"`

	ServiceUsername *string `reform:"service_username"`
	ServicePassword *string `reform:"service_password"`
	ListenPort      *uint16 `reform:"listen_port"`
}

// DSN returns DSN for SNMP service.
func (p *SNMPExporter) DSN(service *SNMPService) string {
	q := make(url.Values)
	q.Set("sslmode", "disable") // TODO https://jira.percona.com/browse/PMM-1727
	q.Set("connect_timeout", strconv.Itoa(int(sqlDialTimeout.Seconds())))

	address := net.JoinHostPort(*service.Address, strconv.Itoa(int(*service.Port)))
	uri := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(*p.ServiceUsername, *p.ServicePassword),
		Host:     address,
		Path:     "postgres",
		RawQuery: q.Encode(),
	}
	return uri.String()
}

// QanAgentsRunOnServer returns qan agents run on server
func QanAgentsRunOnServer(q *reform.Querier) ([]QanAgent, error) {
	rows, err := q.Query(`
SELECT agents.id, agents.type, agents.runs_on_node_id, agents.service_username,
	agents.service_password, agents.listen_port, agents.qan_db_instance_uuid
FROM agents
JOIN nodes ON agents.runs_on_node_id = nodes.id
WHERE nodes.type = ? AND agents.type = ?
`, PMMServerNodeType, QanAgentAgentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := make([]QanAgent, 0)
	for rows.Next() {
		var id, runsOnNodeID int32
		var listenPort sql.NullInt16
		var t AgentType
		var serviceUsername, servicePassword, instanceUUID sql.NullString

		err = rows.Scan(&id, &t, &runsOnNodeID, &serviceUsername, &servicePassword, &listenPort, &instanceUUID)
		if err != nil {
			return nil, err
		}

		agent := QanAgent{
			ID:           id,
			Type:         t,
			RunsOnNodeID: runsOnNodeID,
		}
		if serviceUsername.Valid {
			agent.ServiceUsername = &servicePassword.String
		}
		if servicePassword.Valid {
			agent.ServicePassword = &servicePassword.String
		}
		if instanceUUID.Valid {
			agent.QANDBInstanceUUID = &instanceUUID.String
		}
		if listenPort.Valid {
			port := uint16(listenPort.Int16)
			agent.ListenPort = &port
		}

		agents = append(agents, agent)
	}

	return agents, nil
}
