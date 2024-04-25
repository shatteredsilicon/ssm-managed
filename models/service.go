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
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

//go:generate reform

type ServiceType string

const (
	RDSServiceType        ServiceType = "rds"
	PostgreSQLServiceType ServiceType = "postgresql"
	MySQLServiceType      ServiceType = "mysql"
	SNMPServiceType       ServiceType = "snmp"
)

//reform:services
type Service struct {
	ID     int32       `reform:"id,pk"`
	Type   ServiceType `reform:"type"`
	NodeID int32       `reform:"node_id"`
}

//reform:services
type RDSService struct {
	ID     int32       `reform:"id,pk"`
	Type   ServiceType `reform:"type"`
	NodeID int32       `reform:"node_id"`

	AWSAccessKey  *string `reform:"aws_access_key"` // may be nil
	AWSSecretKey  *string `reform:"aws_secret_key"` // may be nil
	Address       *string `reform:"address"`
	Port          *uint16 `reform:"port"`
	Engine        *string `reform:"engine"`
	EngineVersion *string `reform:"engine_version"`
}

//reform:services
type RDSServiceDetail struct {
	RDSService
	Region   string `reform:"region"`
	Instance string `reform:"instance"`
}

//reform:services
type PostgreSQLService struct {
	ID     int32       `reform:"id,pk"`
	Type   ServiceType `reform:"type"`
	NodeID int32       `reform:"node_id"`

	Address       *string `reform:"address"`
	Port          *uint16 `reform:"port"`
	Engine        *string `reform:"engine"`
	EngineVersion *string `reform:"engine_version"`
}

//reform:services
type MySQLService struct {
	ID     int32       `reform:"id,pk"`
	Type   ServiceType `reform:"type"`
	NodeID int32       `reform:"node_id"`

	Address       *string `reform:"address"`
	Port          *uint16 `reform:"port"`
	Engine        *string `reform:"engine"`
	EngineVersion *string `reform:"engine_version"`
}

//reform:services
type SNMPService struct {
	ID     int32       `reform:"id,pk"`
	Type   ServiceType `reform:"type"`
	NodeID int32       `reform:"node_id"`

	PrivPassword *string `reform:"aws_secret_key"`
	Address      *string `reform:"address"`
	Port         *uint16 `reform:"port"`
	Engine       *string `reform:"engine"`
	// EngineVersion is in format:
	// for v1 or v2 - <version>:<community>
	// for v3 - <version>:<security_level>:<auth_protocol>:<priv_protocol>:<context_name>
	// e.g. 1:public, 2:public, 3:authPriv:MD5:DES:public
	EngineVersion *SNMPEngineVersion `reform:"engine_version"`
}

type SNMPEngineVersion struct {
	Version       string
	Community     string
	SecurityLevel string
	AuthProtocol  string
	PrivProtocol  string
	ContextName   string
}

func (ev *SNMPEngineVersion) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("incompatible type: %s", reflect.ValueOf(src).Kind().String())
	}
	if str == nil {
		return nil
	}

	engineVersion := SNMPEngineVersionFromString(string(str))
	ev.Version = engineVersion.Version
	ev.Community = engineVersion.Community
	ev.SecurityLevel = engineVersion.SecurityLevel
	ev.AuthProtocol = engineVersion.AuthProtocol
	ev.PrivProtocol = engineVersion.PrivProtocol
	ev.ContextName = engineVersion.ContextName

	return nil
}

func (ev SNMPEngineVersion) Value() (driver.Value, error) {
	switch ev.Version {
	case "1", "2":
		return fmt.Sprintf("%s:%s", ev.Version, ev.Community), nil
	default:
		return fmt.Sprintf("%s:%s:%s:%s:%s", ev.Version, ev.SecurityLevel, ev.AuthProtocol, ev.PrivProtocol, ev.ContextName), nil
	}
}

func SNMPEngineVersionFromString(str string) SNMPEngineVersion {
	ev := SNMPEngineVersion{}

	strs := strings.SplitN(str, ":", 2)
	ev.Version = strs[0]
	switch ev.Version {
	case "1", "2":
		if len(strs) > 1 {
			ev.Community = strs[1]
		}
	default:
		strs := strings.SplitN(str, ":", 5)
		strsLen := len(strs)
		if strsLen > 1 {
			ev.SecurityLevel = strs[1]
		}
		if strsLen > 2 {
			ev.AuthProtocol = strs[2]
		}
		if strsLen > 3 {
			ev.PrivProtocol = strs[3]
		}
		if strsLen > 4 {
			ev.ContextName = strs[4]
		}
	}

	return ev
}

//reform:services
type RemoteService struct {
	ID     int32       `reform:"id,pk"`
	Type   ServiceType `reform:"type"`
	NodeID int32       `reform:"node_id"`

	Address       *string `reform:"address"`
	Port          *uint16 `reform:"port"`
	Engine        *string `reform:"engine"`
	EngineVersion *string `reform:"engine_version"`
}

//reform:services
type FullService struct {
	ID     int32       `reform:"id,pk"`
	Type   ServiceType `reform:"type"`
	NodeID int32       `reform:"node_id"`

	AWSAccessKey  *string `reform:"aws_access_key"` // may be nil
	AWSSecretKey  *string `reform:"aws_secret_key"` // may be nil
	Address       *string `reform:"address"`
	Port          *uint16 `reform:"port"`
	Engine        *string `reform:"engine"`
	EngineVersion *string `reform:"engine_version"`
}
