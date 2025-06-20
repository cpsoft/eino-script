package engine

import (
	"eino-script/engine/types"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

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
		case types.NodeTypeChatModel:
			node, err = e.CreateChatModelNode(&cfg)
			break
		case types.NodeTypeChatTemplate:
			node, err = e.CreateChatTemplateNode(&cfg)
			break
		case types.NodeTypeMcpTemplate:
			node, err = e.CreateMcpTemplateNode(&cfg)
			break
		case types.NodeTypeMcpTool:
			node, err = e.CreateMcpToolNode(&cfg)
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
