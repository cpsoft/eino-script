package engine

import (
	"eino-script/engine/nodes"
	"eino-script/engine/types"
)

type ChatTemplateNode struct {
	types.Node
}

func (e ChatTemplateNode) Id() string {
	return e.NodeId
}

func (e ChatTemplateNode) Type() (types.NodeType, error) {
	return e.NodeType, nil
}

func (cm *ChatTemplateNode) GetTargetId() (string, error) {
	return cm.NodeId, nil
}

func (cm *ChatTemplateNode) GetSourceId() (string, error) {
	return cm.NodeId, nil
}

func (e *Engine) CreateChatTemplateNode(cfg *types.NodeCfg) (types.NodeInterface, error) {
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

	template, err := nodes.CreateChatTemplateNode(cfg)
	if err != nil {
		return nil, err
	}

	err = e.g.AddChatTemplateNode(id, template)
	if err != nil {
		return nil, err
	}

	return node, nil
}
