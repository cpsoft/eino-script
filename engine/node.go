package engine

import (
	"context"
	"eino-script/engine/nodes"
	"eino-script/types"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"github.com/sirupsen/logrus"
)

func (e *Engine) CreateEmbeddingNode(cfg *types.NodeCfg) error {
	return nil
}

//func (e *Engine) CreateChatModelNode(cfg *types.NodeCfg) error {
//	data, err := nodes.CreateChatModelByFactroy(cfg)
//	if err != nil {
//		return err
//	}
//	m := data.(model.ToolCallingChatModel)
//
//	bindTool, ok := cfg.Attrs["bindtool"].(string)
//	if ok {
//		toolsInfo, ok := e.tools[bindTool]
//		if !ok {
//			return errors.New("bind tool not found: " + bindTool)
//		}
//		m, err = m.WithTools(toolsInfo)
//		if err != nil {
//			return err
//		}
//	}
//
//	e.models[cfg.Name] = m
//	return e.g.AddChatModelNode(cfg.Name, m)
//}

func (e *Engine) CreateChatTemplateNode(cfg *types.NodeCfg) error {
	data, err := nodes.CreateChatTemplateByFactroy(cfg)
	if err != nil {
		return err
	}
	return e.g.AddChatTemplateNode(cfg.Name, data)
}

// MCPTools太过复杂，进行特殊处理
func (e *Engine) CreateMcpToolsNode(cfg *types.NodeCfg) error {
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

func (e *Engine) CreateToolsNode(cfg *types.NodeCfg) error {
	if cfg.Type == "Mcp" {
		return e.CreateMcpToolsNode(cfg)
	} else {
		//data, err := nodes.CreateNodeByFactroy(cfg)
		//if err != nil {
		//	return err
		//}
		//return e.g.AddToolsNode(cfg.Name, data.(*compose.ToolsNode))
		return nil
	}
}

func (e *Engine) CreateTools(cfgs *[]types.NodeCfg) error {
	for _, cfg := range *cfgs {
		err := e.CreateToolsNode(&cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) CreateChatTemplates(cfgs *[]types.NodeCfg) error {
	for _, cfg := range *cfgs {
		err := e.CreateChatTemplateNode(&cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) CreateNodes(cfgs *[]types.NodeCfg) error {
	var node types.NodeInterface
	var err error

	for _, cfg := range *cfgs {
		switch cfg.Type {
		case types.NodeTypeStart:
			node, err = CreateStartNode(&cfg)
			break
		case types.NodeTypeEnd:
			node, err = CreateEndNode(&cfg)
			break
		case types.NodeTypeChartModel:
			node, err = e.CreateChatModelNode(&cfg)
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
