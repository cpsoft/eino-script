package engine

import (
	"context"
	components "eino-script/components/tools"
	"eino-script/parser"
	"errors"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

func (e *Engine) CreateMcpToolNode(cfg *parser.NodeCfg) error {
	logrus.Infof("CreateMcpToolNode: %+v", *cfg)
	serverName, ok := cfg.Attrs["server"].(string)
	if serverName == "" {
		return errors.New("CreateMcpTemplateNode: name is required")
	}

	mcpServer, ok := e.mcps[serverName]
	if !ok {
		return errors.New("CreateMcpTemplateNode: mcp server not found")
	}

	ctx := context.Background()
	tools, err := components.GetMcpTools(ctx, &components.McpToolsConfig{
		Server:       mcpServer,
		ToolNameList: []string{},
	})
	if err != nil {
		return err
	}

	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})

	toolsInfo := make([]*schema.ToolInfo, 0)
	for _, tool := range tools {
		info, err := tool.Info(ctx)
		if err == nil {
			toolsInfo = append(toolsInfo, info)
		}

	}
	e.tools[cfg.Name] = toolsInfo
	_ = e.g.AddToolsNode(cfg.Name, toolsNode)
	return nil
}
