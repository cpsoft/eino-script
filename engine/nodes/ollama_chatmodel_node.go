package nodes

import (
	"context"
	"eino-script/types"
	"errors"
	"fmt"
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

	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in attrs")
	}

	BaseUrl, ok := data["base_url"].(string)
	if !ok {
		//return nil, errors.New("base url not found in config")
		BaseUrl = "http://localhost:11434/"
	}

	Model, ok := data["model"].(string)
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
