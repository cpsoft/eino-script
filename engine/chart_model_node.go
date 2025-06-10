package engine

import (
	"eino-script/engine/nodes"
	"eino-script/types"
)

type ChatModelNode struct {
	types.Node
	templateId  string
	chatModelId string
}

func (e ChatModelNode) Id() string {
	return e.NodeId
}

func (e ChatModelNode) Type() (types.NodeType, error) {
	return e.NodeType, nil
}

func (cm *ChatModelNode) GetTargetId() (string, error) {
	return cm.templateId, nil
}

func (cm *ChatModelNode) GetSourceId() (string, error) {
	return cm.chatModelId, nil
}

func (e *Engine) CreateChatModelNode(cfg *types.NodeCfg) (types.NodeInterface, error) {
	n, err := CreateGeneralNode(cfg)
	if err != nil {
		return nil, err
	}

	node := &ChatModelNode{
		Node: *n,
	}

	id := node.Id()
	if id == "" {
		return nil, err
	}

	node.templateId = id + "-" + "template"
	template, err := nodes.CreateChatTemplateNode(cfg)
	if err != nil {
		return nil, err
	}

	err = e.g.AddChatTemplateNode(node.templateId, template)
	if err != nil {
		return nil, err
	}

	node.chatModelId = id + "-" + "chatmodel"
	chatModel, err := nodes.CreateOllamaChatModelNode(cfg)
	if err != nil {
		return nil, err
	}
	err = e.g.AddChatModelNode(node.chatModelId, chatModel)
	if err != nil {
		return nil, err
	}

	err = e.g.AddEdge(node.templateId, node.chatModelId)
	if err != nil {
		return nil, err
	}
	return node, nil
}
