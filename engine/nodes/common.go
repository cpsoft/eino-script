package nodes

import (
	"eino-script/types"
	"errors"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
)

type NodesFactroyData struct {
}

var chattemplate_functions = map[types.NodeType]func(cfg *types.NodeCfg) (prompt.ChatTemplate, error){
	"ChatTemplate": CreateChatTemplateNode,
}

func CreateChatTemplateByFactroy(cfg *types.NodeCfg) (prompt.ChatTemplate, error) {
	f, ok := chattemplate_functions[cfg.Type]
	if !ok {
		return nil, errors.New("Node not found")
	}
	return f(cfg)
}

var chatmodel_functions = map[types.NodeType]func(cfg *types.NodeCfg) (model.ToolCallingChatModel, error){
	"Ollama": CreateOllamaChatModelNode,
	"Qwen":   CreateQwenChatModelNode,
}

func CreateChatModelByFactroy(cfg *types.NodeCfg) (model.ToolCallingChatModel, error) {
	f, ok := chatmodel_functions[cfg.Type]
	if !ok {
		return nil, errors.New("Node not found")
	}
	return f(cfg)
}
