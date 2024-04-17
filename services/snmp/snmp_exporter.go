package snmp

import (
	"github.com/pkg/errors"
	"github.com/shatteredsilicon/snmp_exporter/config"
	"gopkg.in/yaml.v2"
)

type snmpGeneratorConfig struct {
	Auths   map[string]*config.Auth         `yaml:"auths"`
	Modules map[string]*snmpGeneratorModule `yaml:"modules"`
	Version int                             `yaml:"version,omitempty"`
}

type snmpGeneratorModule struct {
	Walk       []string                          `yaml:"walk"`
	Lookups    []*snmpGeneratorLookup            `yaml:"lookups"`
	WalkParams config.WalkParams                 `yaml:",inline"`
	Overrides  map[string]snmpGeneratorOverrides `yaml:"overrides"`
	Filters    config.Filters                    `yaml:"filters,omitempty"`
}

type snmpGeneratorLookup struct {
	SourceIndexes     []string `yaml:"source_indexes"`
	Lookup            string   `yaml:"lookup"`
	DropSourceIndexes bool     `yaml:"drop_source_indexes,omitempty"`
}

type snmpGeneratorOverrides struct {
	Ignore         bool                              `yaml:"ignore,omitempty"`
	RegexpExtracts map[string][]config.RegexpExtract `yaml:"regex_extracts,omitempty"`
	Offset         float64                           `yaml:"offset,omitempty"`
	Scale          float64                           `yaml:"scale,omitempty"`
	Type           string                            `yaml:"type,omitempty"`
}

func (config *snmpGeneratorConfig) Marshal() ([]byte, error) {
	b, err := yaml.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal snmp_exporter generator configuration")
	}
	b = append([]byte("# Managed by ssm-managed. DO NOT EDIT.\n"), b...)
	return b, nil
}
