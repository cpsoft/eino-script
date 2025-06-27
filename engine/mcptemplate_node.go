package engine

import (
	"eino-script/engine/templates"
	"eino-script/engine/types"
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

func (e *Engine) CreateMcpTemplateNode(cfg *types.NodeCfg) (types.NodeInterface, error) {
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

	err = e.g.AddChatTemplateNode(id, template)
	if err != nil {
		return nil, err
	}

	return node, nil
}
