package nodes

import (
	"context"
	"eino-script/types"
	"errors"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/sirupsen/logrus"
)

func CreateQwenChatModelNode(cfg *types.NodeCfg) (model.ToolCallingChatModel, error) {
	logrus.Infof("CreateOllamaChatModelNode: %+v", *cfg)
	BaseUrl, ok := cfg.Attrs["base_url"].(string)
	if !ok {
		return nil, errors.New("base url not found in config")
	}
	ApiKey, ok := cfg.Attrs["api_key"].(string)
	if !ok {
		return nil, errors.New("api key not found in config")
	}

	Model, ok := cfg.Attrs["model"].(string)
	if !ok {
		return nil, errors.New("model not found in config")
	}

	model, err := qwen.NewChatModel(context.Background(), &qwen.ChatModelConfig{
		BaseURL: BaseUrl,
		APIKey:  ApiKey,
		Model:   Model,
		Timeout: 0,
	})
	if err != nil {
		return nil, err
	}

	return model, nil
}
