package nodes

import (
	"context"
	types2 "eino-script/engine/types"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/sirupsen/logrus"
)

func CreateOpenaiChatModelNode(info *types2.ModelInfo, cfg *types2.NodeCfg) (model.ToolCallingChatModel, error) {
	logrus.Infof("CreateOpenaiChatModelNode: %+v", *cfg)

	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in attrs")
	}

	Temperature, ok := data["temperature"].(float32)
	if !ok {
		Temperature = 0.7
	}

	model, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL:     info.ApiUrl,
		Model:       info.ModelName,
		APIKey:      info.ApiKey,
		Temperature: &Temperature,
	})
	if err != nil {
		return nil, err
	}

	return model, nil
}
