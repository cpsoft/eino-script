package engine

import (
	"eino-script/engine/components"
	"eino-script/engine/types"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

func (e *Engine) CreateNodes(cfgs *[]types.NodeCfg) error {
	var node types.NodeInterface
	var err error

	for _, cfg := range *cfgs {
		node = nil
		switch cfg.Type {
		case types.NodeTypeStart:
			node, err = components.CreateStartNode(&cfg)
			break
		case types.NodeTypeEnd:
			node, err = components.CreateEndNode(&cfg)
			break
		case types.NodeTypeChatModel:
			node, err = components.CreateChatModelNode(&cfg, e.g, e.callbacks)
			break
		case types.NodeTypeChatTemplate:
			node, err = components.CreateChatTemplateNode(&cfg, e.g)
			break
		case types.NodeTypeMcpTemplate:
			node, err = components.CreateMcpTemplateNode(&cfg, e.g)
			break
		case types.NodeTypeMcpTool:
			node, err = components.CreateMcpToolNode(&cfg, e.g, e.callbacks)
			break
		case types.NodeTypeLoader:
			node, err = components.CreateLoaderNode(&cfg, e.g)
			break
		case types.NodeTypeBranch:
			branch, err := components.CreateBranch(&cfg)
			if err != nil {
				return fmt.Errorf("创建Branch失败(%s)：%s", cfg.Id, err.Error())
			}
			e.branchs[branch.Id()] = branch
			break
		case types.NodeTypeLoop:
			break
		default:
			err = errors.New(string("Unknown node type:" + cfg.Type))
			break
		}
		if err != nil {
			return fmt.Errorf("创建节点失败(%s)：%s", cfg.Id, err.Error())
		}

		if node != nil {
			e.nodes[node.Id()] = node
		}
	}
	logrus.Debug(e.nodes)
	return nil
}
