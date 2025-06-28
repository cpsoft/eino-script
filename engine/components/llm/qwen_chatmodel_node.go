package llm

import (
	"context"
	types2 "eino-script/engine/types"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/sirupsen/logrus"
)

func CreateQwenChatModelNode(info types2.ModelInfo, cfg *types2.NodeCfg) (model.ToolCallingChatModel, error) {
	logrus.Infof("CreateOllamaChatModelNode: %+v", *cfg)
	model, err := qwen.NewChatModel(context.Background(), &qwen.ChatModelConfig{
		BaseURL: info.ApiUrl,
		APIKey:  info.ApiKey,
		Model:   info.ModelName,
		Timeout: 0,
	})
	if err != nil {
		return nil, err
	}

	return model, nil
}
