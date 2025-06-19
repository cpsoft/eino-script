package engine

import (
	"eino-script/engine/nodes"
	types2 "eino-script/engine/types"
)

type McpTemplateNode struct {
	types2.Node
}

func (e McpTemplateNode) Id() string {
	return e.NodeId
}

func (e McpTemplateNode) Type() (types2.NodeType, error) {
	return e.NodeType, nil
}

func (cm *McpTemplateNode) GetTargetId() (string, error) {
	return cm.NodeId, nil
}

func (cm *McpTemplateNode) GetSourceId() (string, error) {
	return cm.NodeId, nil
}

func (e *Engine) CreateMcpTemplateNode(cfg *types2.NodeCfg) (types2.NodeInterface, error) {
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

	template, err := nodes.CreateMcpTemplateNode(cfg)
	if err != nil {
		return nil, err
	}

	err = e.g.AddChatTemplateNode(id, template)
	if err != nil {
		return nil, err
	}

	return node, nil
}
