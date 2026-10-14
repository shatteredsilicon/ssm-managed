package agent

import (
	"context"
	"fmt"
	"os"
	"path"

	servicelib "github.com/percona/kardianos-service"
	"github.com/pkg/errors"
	"github.com/shatteredsilicon/ssm-managed/models"
	"github.com/shatteredsilicon/ssm-managed/services"
	"gopkg.in/ini.v1"
	"gopkg.in/reform.v1"
)

type ToolAgent struct {
	toolAgentPath   string
	baseDir         string
	ssmServerNodeID int32
	supervisor      services.Supervisor
}

type ToolAgentConfItem struct {
	Name       string `json:"name"`
	Executable string `json:"executable"`
}

func NewToolAgent(
	ssmServerNodeID int32,
	toolAgentPath, baseDir string,
	supervisor services.Supervisor,
) *ToolAgent {
	return &ToolAgent{
		toolAgentPath:   toolAgentPath,
		baseDir:         baseDir,
		ssmServerNodeID: ssmServerNodeID,
		supervisor:      supervisor,
	}
}

func (ta *ToolAgent) Add(ctx context.Context, tx *reform.TX, service *models.Service, instance string, confItems []ToolAgentConfItem) error {
	toolAgentCfgPath := path.Join(ta.baseDir, toolAgentConfName(instance))
	toolAgentCfgFile, err := ini.LooseLoad(toolAgentCfgPath)
	if err != nil {
		return err
	}

	logfileSuffix := "-" + instance
	if instance == models.SSMServerNodeName {
		logfileSuffix = ""
	}
	for _, item := range confItems {
		toolAgentCfgFile.Section(item.Name).Key("executable").SetValue(item.Executable)
		toolAgentCfgFile.Section(item.Name).Key("output").SetValue(path.Join(ta.baseDir, fmt.Sprintf("%s%s.log", item.Name, logfileSuffix)))
	}

	if err := toolAgentCfgFile.SaveTo(toolAgentCfgPath); err != nil {
		return err
	}

	serviceName := ta.ToolAgentServiceName(instance)
	if instance == models.SSMServerNodeName {
		if err := ta.supervisor.Status(ctx, serviceName); err == nil {
			if err = ta.supervisor.StopService(ctx, serviceName); err != nil {
				return err
			}
		}

		return ta.supervisor.StartService(ctx, serviceName)
	} else {
		if tx != nil { // tx is not nil means it's in a progress that needs to create the agent record
			agent := &models.Agent{
				Type:         models.MySQLdExporterAgentType,
				RunsOnNodeID: ta.ssmServerNodeID,
			}
			if err = tx.Insert(agent); err != nil {
				return errors.WithStack(err)
			}
			if err = tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID}); err != nil {
				return errors.WithStack(err)
			}
		}

		if err := ta.supervisor.Status(ctx, serviceName); err == nil {
			if err = ta.supervisor.Stop(ctx, serviceName); err != nil {
				return err
			}
		}

		cfg := ta.toolAgentCfg(instance, toolAgentCfgPath)
		return ta.supervisor.Start(ctx, cfg)
	}
}

func (ta *ToolAgent) Remove(ctx context.Context, instance string) error {
	serviceName := ta.ToolAgentServiceName(instance)
	toolAgentCfgPath := path.Join(ta.baseDir, toolAgentConfName(instance))

	if serviceName == models.SSMServerNodeName {
		toolAgentCfgFile, err := ini.LooseLoad(toolAgentCfgPath)
		if err != nil {
			return err
		}

		if err = ta.supervisor.StopService(ctx, serviceName); err != nil {
			return err
		}

		for _, sectionName := range toolAgentCfgFile.SectionStrings() {
			if sectionName == ini.DefaultSection || sectionName == "" {
				continue
			}

			toolAgentCfgFile.DeleteSection(sectionName)
		}

		if err := toolAgentCfgFile.SaveTo(toolAgentCfgPath); err != nil {
			return err
		}

		if err := ta.supervisor.StartService(ctx, serviceName); err != nil {
			return err
		}
	} else {
		if err := ta.supervisor.Stop(ctx, serviceName); err != nil {
			return err
		}
		if err := os.Remove(toolAgentCfgPath); err != nil {
			return err
		}
	}

	return nil
}

func (ta *ToolAgent) toolAgentCfg(instance, toolAgentCfgPath string) *servicelib.Config {
	name := toolAgentConfName(instance)

	arguments := []string{
		fmt.Sprintf("-config=\"%s\"", toolAgentCfgPath),
	}

	return &servicelib.Config{
		Name:        name,
		DisplayName: name,
		Description: name,
		Executable:  ta.toolAgentPath,
		Arguments:   arguments,
	}
}

func (ta *ToolAgent) ToolAgentServiceName(instance string) string {
	if instance == models.SSMServerNodeName {
		return string(models.ToolAgentAgentType)
	}
	return fmt.Sprintf("ssm-tool-agent-%s", instance)
}

func toolAgentConfName(instance string) string {
	if instance == models.SSMServerNodeName {
		return fmt.Sprintf("%s.conf", string(models.ToolAgentAgentType))
	}
	return fmt.Sprintf("ssm-%s-%s.conf", string(models.ToolAgentAgentType), instance)
}
