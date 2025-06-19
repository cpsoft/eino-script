package engine

import (
	"context"
	"eino-script/engine/nodes"
	types2 "eino-script/engine/types"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"github.com/sirupsen/logrus"
)

// MCPTools太过复杂，进行特殊处理
func (e *Engine) CreateMcpToolsNode(cfg *types2.NodeCfg) error {
	serverName, ok := cfg.Attrs["server"].(string)
	if serverName == "" {
		return errors.New("CreateMcpTemplateNode: name is required")
	}

	mcpServer, ok := e.mcps[serverName]
	if !ok {
		return errors.New("CreateMcpTemplateNode: mcp server not found")
	}

	toolsInfo, err := mcpServer.ListTools(context.Background(), []string{})
	if err != nil {
		return err
	}
	e.tools[cfg.Name] = toolsInfo
	data, err := nodes.CreateMcpToolNode(cfg, mcpServer, toolsInfo)
	if err != nil {
		return err
	}
	return e.g.AddToolsNode(cfg.Name, data.(*compose.ToolsNode))
}

func (e *Engine) CreateNodes(cfgs *[]types2.NodeCfg) error {
	var node types2.NodeInterface
	var err error

	for _, cfg := range *cfgs {
		switch cfg.Type {
		case types2.NodeTypeStart:
			node, err = CreateStartNode(&cfg)
			break
		case types2.NodeTypeEnd:
			node, err = CreateEndNode(&cfg)
			break
		case types2.NodeTypeChatModel:
			node, err = e.CreateChatModelNode(&cfg)
			break
		case types2.NodeTypeChatTemplate:
			node, err = e.CreateChatTemplateNode(&cfg)
			break
		case types2.NodeTypeMcpTemplate:
			node, err = e.CreateMcpTemplateNode(&cfg)
			break
		default:
			err = errors.New(string("Unknown node type:" + cfg.Type))
			break
		}
		if err != nil {
			return fmt.Errorf("创建节点失败(%s)：%s", cfg.Id, err.Error())
		}

		e.nodes[node.Id()] = node
	}
	logrus.Debug(e.nodes)
	return nil
}
