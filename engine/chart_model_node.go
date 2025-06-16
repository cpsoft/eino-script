package engine

import (
	"eino-script/engine/nodes"
	"eino-script/types"
	"fmt"
	"github.com/cloudwego/eino/components/model"
	"github.com/sirupsen/logrus"
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

	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in attrs")
	}
	logrus.Debug(data)

	ModelId, ok := data["model"].(string)
	if !ok {
		return nil, fmt.Errorf("model not found in config")
	}

	if e.callbacks == nil {
		return nil, fmt.Errorf("engine的回调函数没有配置。")
	}
	info, err := e.callbacks.GetModelInfo(ModelId)
	if err != nil {
		return nil, err
	}

	var chatModel model.ToolCallingChatModel
	switch info.ModelType {
	case "ollama":
		chatModel, err = nodes.CreateOllamaChatModelNode(info, cfg)
	case "openai":
		chatModel, err = nodes.CreateOpenaiChatModelNode(info, cfg)
	default:
		return nil, fmt.Errorf("模型类型不正确:(" + info.ModelType + ")")
	}

	if err != nil {
		return nil, err
	}

	node.chatModelId = id + "-" + "chatmodel"
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
