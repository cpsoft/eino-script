package nodes

import (
	"context"
	components "eino-script/components/tools"
	"eino-script/types"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

func CreateMcpToolNode(cfg *types.NodeCfg, server types.IMcpServer, toolsInfo []*schema.ToolInfo) (interface{}, error) {
	logrus.Infof("CreateMcpToolNode: %+v", *cfg)
	ctx := context.Background()
	tools, err := components.GetMcpTools(ctx, &components.McpToolsConfig{
		Server:    server,
		ToolsInfo: toolsInfo,
	})
	if err != nil {
		return nil, err
	}

	return compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})
}
