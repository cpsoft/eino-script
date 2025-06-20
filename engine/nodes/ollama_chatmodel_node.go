package nodes

import (
	"context"
	"eino-script/engine/types"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
)

func CreateOllamaChatModelNode(
	info *types.ModelInfo,
	cfg *types.NodeCfg,
	tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	var err error
	logrus.Infof("CreateOllamaChatModelNode: %+v", *cfg)

	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in attrs")
	}

	Temperature, ok := data["temperature"].(float32)
	if !ok {
		Temperature = 0.7
	}

	var model model.ToolCallingChatModel
	model, err = ollama.NewChatModel(context.Background(), &ollama.ChatModelConfig{
		BaseURL: info.ApiUrl,
		Model:   info.ModelName,
		Options: &api.Options{
			Temperature: Temperature,
		},
	})
	if err != nil {
		return nil, err
	}

	if tools != nil {
		model, err = model.WithTools(tools)
		if err != nil {
			return nil, err
		}
	}

	return model, nil
}
