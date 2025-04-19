package engine

import (
	"context"
	"eino-script/parser"
	"errors"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/sirupsen/logrus"
)

func (e *Engine) CreateQwenChatModelNode(cfg *parser.NodeCfg) error {
	logrus.Infof("CreateOllamaChatModelNode: %+v", *cfg)
	BaseUrl, ok := cfg.Attrs["base_url"].(string)
	if !ok {
		return errors.New("base url not found in config")
	}
	ApiKey, ok := cfg.Attrs["api_key"].(string)
	if !ok {
		return errors.New("api key not found in config")
	}

	Model, ok := cfg.Attrs["model"].(string)
	if !ok {
		return errors.New("model not found in config")
	}

	model, err := qwen.NewChatModel(context.Background(), &qwen.ChatModelConfig{
		BaseURL: BaseUrl,
		APIKey:  ApiKey,
		Model:   Model,
		Timeout: 0,
	})
	if err != nil {
		return err
	}

	e.models[cfg.Name] = model
	_ = e.g.AddChatModelNode(cfg.Name, model)
	return nil
}
