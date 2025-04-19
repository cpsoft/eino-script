package engine

import (
	"eino-script/parser"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/ollama/ollama/api"
)

type OllamaConfig struct {
	BaseUrl string
	Model   string
}

func (e *Engine) CreateOllamaChatModelNode(cfg *parser.NodeCfg) error {
	fmt.Println("CreateOllamaChatModelNode: ", *cfg)
	BaseUrl, ok := cfg.Attrs["base_url"].(string)
	if !ok {
		return errors.New("base url not found in config")
	}

	Model, ok := cfg.Attrs["model"].(string)
	if !ok {
		return errors.New("model not found in config")
	}

	model, err := ollama.NewChatModel(e.ctx, &ollama.ChatModelConfig{
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
		return err
	}

	e.models[cfg.Name] = model
	_ = e.g.AddChatModelNode(cfg.Name, model)
	return nil
}
