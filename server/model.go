package server

import (
	"database/sql"
	"eino-script/engine"
	"eino-script/types"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

// 查询 ID 和 Name 列表
func fetchModelItems(db *sql.DB) ([]types.ModelInfo, error) {
	query := `SELECT id, name, modelType, apiKey, apiUrl, maxContextLength, streamingEnabled FROM models`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	items := make([]types.ModelInfo, 0)
	for rows.Next() {
		var item types.ModelInfo
		if err := rows.Scan(
			&item.ID,
			&item.ModelName,
			&item.ModelType,
			&item.ApiKey,
			&item.ApiUrl,
			&item.MaxContextLength,
			&item.StreamingEnabled,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	return items, nil
}

func (s *Server) handleGetModelList(c *gin.Context) {
	// 查询数据
	items, err := fetchModelItems(s.db)
	logrus.Debug("items:", items)
	if err != nil {
		Error(c, 300, "读取数据失败："+err.Error())
		return
	}

	Success(c, items)

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

	// 插入或更新数据到数据库
	upsertSQL := `
	INSERT INTO models (id, name, modelType, apiKey, apiUrl, maxContextLength, streamingEnabled) 
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET 
		name = excluded.name, 
		modelType = excluded.modelType,
	           apiKey = excluded.apiKey,
	           apiUrl = excluded.apiUrl,
	           maxContextLength = excluded.maxContextLength,
	           streamingEnabled = excluded.streamingEnabled
	           
	`
	_, err = s.db.Exec(upsertSQL,
		model.ID, model.ModelName,
		model.ModelType, model.ApiKey, model.ApiUrl,
		model.MaxContextLength, model.StreamingEnabled)
	if err != nil {
		logrus.Debug("大模型插入错误：", err.Error())
		Error(c, 200, "插入大模型错误："+err.Error())
		return
	}

	Success(c, gin.H{"message": "保存成功"})
}

type DeleteModelRequest struct {
	ID string `json:"id"`
}

// saveFlow 函数：保存 JSON 数据到 SQLite
func (s *Server) deleteModel(id string) error {
	logrus.Debug("删除数据：", id)

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM models WHERE id = ?", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check flow existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("flow with id %s not found", id)
	}

	deleteSQL := `DELETE FROM models WHERE id=?`
	result, err := s.db.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item with id %s not found", id)
	}

	return nil
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
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	logrus.Debug("delete id:", body.ID)
	err = s.deleteModel(body.ID)
	if err != nil {
		logrus.Debug("删除失败：" + err.Error())
		Error(c, 100, "删除失败："+err.Error())
		return
	}

	logrus.Debug("删除数据成功")
	Success(c, "")
}

type GetOllamaModelNamesRequest struct {
	Url string `json:"url"`
}

func (s *Server) handleOllamaModelNames(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var body GetOllamaModelNamesRequest
	err = json.Unmarshal(data, &body)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}
	baseUrl, err := url.Parse(body.Url)
	if err != nil {
		Error(c, http.StatusBadRequest, "Ollama 服务URL错误。"+body.Url)
		return
	}

	models, err := engine.GetOllamaModels(baseUrl)
	if err != nil {
		Error(c, http.StatusInternalServerError, "获取模型失败："+err.Error())
		return
	}

	Success(c, models)
}

// 根据 ID 获取模型记录
func (s *Server) getModelByID(id string) (*types.ModelInfo, error) {
	// 准备 SQL 查询语句
	query := "SELECT id, name, modelType, apiKey, apiUrl, maxContextLength, StreamingEnabled FROM models WHERE id = ?"
	row := s.db.QueryRow(query, id)

	// 解析查询结果
	var model types.ModelInfo
	err := row.Scan(&model.ID,
		&model.ModelName,
		&model.ModelType,
		&model.ApiKey,
		&model.ApiUrl,
		&model.MaxContextLength,
		&model.StreamingEnabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no record found for ID %d", id)
		}
		return nil, err
	}

	return &model, nil
}

func (s *Server) GetModelInfo(id string) (*types.ModelInfo, error) {
	return s.getModelByID(id)
}
