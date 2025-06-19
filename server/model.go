package server

import (
	"eino-script/engine"
	"eino-script/engine/types"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (s *Server) handleGetModelList(c *gin.Context) {
	// 查询数据
	models, err := s.provider.GetModelList()
	if err != nil {
		Error(c, 300, "读取数据失败："+err.Error())
		return
	}

	Success(c, models)
}

func (s *Server) handleSaveModel(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debug(string(body))

	var model types.ModelInfo
	err = json.Unmarshal(body, &model)
	if err != nil {
		logrus.Error("模型数据解析错误:", err.Error())
		Error(c, http.StatusBadRequest, "模型数据解析错误")
		return
	}

	err = s.provider.SaveModel(&model)
	if err != nil {
		logrus.Debug("大模型插入错误：", err.Error())
		Error(c, 200, err.Error())
		return
	}
	Success(c, gin.H{"message": "保存成功"})
}

type DeleteModelRequest struct {
	ID uint `json:"id"`
}

func (s *Server) handleDeleteModel(c *gin.Context) {
	data, err := c.GetRawData()
	logrus.Debug("delete data:", string(data))
	if err != nil {
		Error(c, http.StatusBadRequest, "数据请求错误")
		return
	}

	var body DeleteModelRequest
	err = json.Unmarshal(data, &body)
	if err != nil {
		Error(c, http.StatusBadRequest, "请求数据格式错误。")
		return
	}

	logrus.Debug("delete id:", body.ID)
	err = s.provider.DeleteModel(body.ID)
	if err != nil {
		logrus.Debug(err.Error())
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Debug("删除数据成功")
	Success(c, "")
}

type GetOllamaModelNamesRequest struct {
	Url string `json:"url"`
}

func (s *Server) handleChatModelList(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var body types.ModelInfo
	err = json.Unmarshal(data, &body)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	models := make([]string, 0)
	modelType := strings.ToLower(strings.TrimSpace(body.ModelType))
	logrus.Debug("modelType:", modelType)
	switch modelType {
	case "ollama":
		models, err = engine.GetOllamaModels(&body)
	case "openai":
		models, err = engine.GetOpenaiModels(&body)
	default:
		err = errors.New("模型类型不正确")
	}

	if err != nil {
		Error(c, http.StatusInternalServerError, "获取模型失败："+err.Error())
		return
	}

	logrus.Debug("list:", models)

	Success(c, models)
}

// Callbacks
func (s *Server) Callback_GetModelInfo(modelId uint) (*types.ModelInfo, error) {
	model, err := s.provider.GetModel(modelId)
	if err != nil {
		return nil, err
	}
	info := types.ModelInfo{
		ID:               model.ID,
		Name:             model.Name,
		ModelType:        model.ModelType,
		ModelName:        model.ModelName,
		ApiKey:           model.ApiKey,
		ApiUrl:           model.ApiUrl,
		MaxContextLength: model.MaxContextLength,
		StreamingEnabled: model.StreamingEnabled,
	}
	return &info, nil
}
