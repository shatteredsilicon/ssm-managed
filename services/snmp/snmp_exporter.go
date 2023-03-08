package snmp

import (
	"github.com/pkg/errors"
	"github.com/shatteredsilicon/snmp_exporter/config"
	"gopkg.in/yaml.v2"
)

type snmpGeneratorConfig struct {
	Modules map[string]*snmpGeneratorModule `yaml:"modules"`
}

type snmpGeneratorModule struct {
	Walk       []string               `yaml:"walk"`
	WalkParams config.WalkParams      `yaml:",inline"`
	Lookups    []interface{}          `yaml:"lookups"`
	Overrides  map[string]interface{} `yaml:"overrides"`
}

func (config *snmpGeneratorConfig) Marshal() ([]byte, error) {
	b, err := yaml.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal snmp_exporter generator configuration")
	}
	b = append([]byte("# Managed by ssm-managed. DO NOT EDIT.\n"), b...)
	return b, nil
}
