package nodes

import (
	"context"
	"eino-script/types"
	"errors"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
)

type OllamaConfig struct {
	BaseUrl string
	Model   string
}

func CreateOllamaChatModelNode(cfg *types.NodeCfg) (model.ToolCallingChatModel, error) {
	logrus.Infof("CreateOllamaChatModelNode: %+v", *cfg)
	BaseUrl, ok := cfg.Attrs["base_url"].(string)
	if !ok {
		return nil, errors.New("base url not found in config")
	}

	Model, ok := cfg.Attrs["model"].(string)
	if !ok {
		return nil, errors.New("model not found in config")
	}

	model, err := ollama.NewChatModel(context.Background(), &ollama.ChatModelConfig{
		BaseURL: BaseUrl,
		Model:   Model,

		//Timeout: 30 * time.Second,
		//Format:  json.RawMessage(`"json"`),
		Options: &api.Options{
			Temperature: 0.7,
			//NumPredict:  100,
		},
	})
	if err != nil {
		return nil, err
	}

	return model, nil
}
