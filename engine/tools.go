package engine

import (
	"context"
	"errors"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
)

func GetOllamaModels(baseUrl *url.URL) ([]string, error) {
	client := api.NewClient(baseUrl, http.DefaultClient)

	models, err := client.List(context.Background())
	if err != nil {
		return nil, errors.New("获取模型列表失败。")
	}

	list := make([]string, 0)
	for _, m := range models.Models {
		list = append(list, m.Name)
	}
	return list, nil
}
