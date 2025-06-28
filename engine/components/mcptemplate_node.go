package components

import (
	"eino-script/engine/components/templates"
	"eino-script/engine/types"
	"github.com/cloudwego/eino/compose"
)

type McpTemplateNode struct {
	types.Node
}

func (e McpTemplateNode) Id() string {
	return e.NodeId
}

func (e McpTemplateNode) Type() (types.NodeType, error) {
	return e.NodeType, nil
}

func (cm *McpTemplateNode) GetTargetId() (string, error) {
	return cm.NodeId, nil
}

func (cm *McpTemplateNode) GetSourceId() (string, error) {
	return cm.NodeId, nil
}

func CreateMcpTemplateNode(cfg *types.NodeCfg, g *compose.Graph[any, any]) (types.NodeInterface, error) {
	n, err := CreateGeneralNode(cfg)
	if err != nil {
		return nil, err
	}

	node := &ChatTemplateNode{
		Node: *n,
	}

	id := node.Id()
	if id == "" {
		return nil, err
	}

	template, err := templates.CreateMcpTemplateNode(cfg)
	if err != nil {
		return nil, err
	}

	err = g.AddChatTemplateNode(id, template)
	if err != nil {
		return nil, err
	}

	return node, nil
}
