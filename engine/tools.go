package engine

import (
	"context"
	"eino-script/engine/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

func GetOllamaModels(info *types.ModelInfo) ([]string, error) {
	baseUrl, err := url.Parse(info.ApiUrl)
	if err != nil {
		return nil, err
	}
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

type OpenaiModel struct {
	ID      string `json:"ID"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type OpenaiModelsResponse struct {
	Object string        `json:"object"`
	Data   []OpenaiModel `json:"data"`
}

func GetOpenaiModels(info *types.ModelInfo) ([]string, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + info.ApiKey,
	}

	req, err := http.NewRequest("GET", info.ApiUrl+"/models", nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			err = fmt.Errorf("ApiKey设置错误。")
			break
		default:
			err = fmt.Errorf("未知错误：", string(body))
			break
		}
		return nil, err
	}

	var response OpenaiModelsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)
	for _, m := range response.Data {
		if m.Object == "model" {
			list = append(list, m.ID)
		} else {
			logrus.Debug("Object: ", m.Object)
		}

	}
	return list, nil
}
