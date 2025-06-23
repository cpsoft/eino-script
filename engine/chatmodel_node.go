package engine

import (
	"context"
	"eino-script/engine/nodes"
	"eino-script/engine/types"
	"fmt"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

type ChatModelNode struct {
	types.Node
}

func (e ChatModelNode) Id() string {
	return e.NodeId
}

func (e ChatModelNode) Type() (types.NodeType, error) {
	return e.NodeType, nil
}

func (cm *ChatModelNode) GetTargetId() (string, error) {
	return cm.NodeId, nil
}

func (cm *ChatModelNode) GetSourceId() (string, error) {
	return cm.NodeId, nil
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

	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in attrs")
	}
	logrus.Debug(data)

	ModelId, ok := data["model"].(float64)
	if !ok {
		return nil, fmt.Errorf("model not found in config")
	}

	if e.callbacks == nil {
		return nil, fmt.Errorf("engine的回调函数没有配置。")
	}
	info, err := e.callbacks.Callback_GetModelInfo(uint(ModelId))
	if err != nil {
		return nil, err
	}

	var chatModel model.ToolCallingChatModel
	var tools []*schema.ToolInfo = nil
	McpId, ok := data["mcpid"].(float64)
	if ok {
		logrus.Debugf("获取到MCP ID: %f", McpId)
		if e.callbacks != nil {
			mcpServer, err := e.callbacks.Callback_CreateMcpServer(uint(McpId))
			if err != nil {
				return nil, fmt.Errorf("大模型创建MCP服务失败")
			}
			tools, err = mcpServer.ListTools(context.Background(), nil)
			if err != nil {
				return nil, err
			}
		}
	} else {
		logrus.Debug("未获取到Mcp ID.")
	}

	switch info.ModelType {
	case "ollama":
		chatModel, err = nodes.CreateOllamaChatModelNode(info, cfg, tools)
	case "openai":
		chatModel, err = nodes.CreateOpenaiChatModelNode(info, cfg, tools)
	default:
		return nil, fmt.Errorf("模型类型不正确:(" + info.ModelType + ")")
	}

	if err != nil {
		return nil, err
	}

	err = e.g.AddChatModelNode(id, chatModel)
	if err != nil {
		return nil, err
	}

	return node, nil
}
