package server

import (
	"context"
	"database/sql"
	"eino-script/components"
	"eino-script/types"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 查询 ID 和 Name 列表
func fetchMcpItems(db *sql.DB) ([]types.McpInfo, error) {
	query := `SELECT id, name, mcpType, url FROM mcps`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	items := make([]types.McpInfo, 0)
	for rows.Next() {
		var item types.McpInfo
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.McpType,
			&item.Url,
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

func (s *Server) handleGetMcpList(c *gin.Context) {
	// 查询数据
	items, err := fetchMcpItems(s.db)
	logrus.Debug("items:", items)
	if err != nil {
		Error(c, 300, "读取数据失败："+err.Error())
		return
	}

	Success(c, items)

}

func (s *Server) handleSaveMcp(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debug(string(body))

	var item types.McpInfo
	err = json.Unmarshal(body, &item)
	if err != nil {
		logrus.Error("MCP数据解析错误:", err.Error())
		Error(c, http.StatusBadRequest, "MCP数据解析错误")
		return
	}

	// 插入或更新数据到数据库
	upsertSQL := `
	INSERT INTO mcps (id, name, mcpType, url) 
	VALUES (?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET 
		name = excluded.name, 
		mcpType = excluded.mcpType,
	    url = excluded.url
	`
	_, err = s.db.Exec(upsertSQL,
		item.ID, item.Name,
		item.McpType, item.Url)
	if err != nil {
		logrus.Debug("MCP插入错误：", err.Error())
		Error(c, 200, "插入MCP错误："+err.Error())
		return
	}

	Success(c, gin.H{"message": "保存成功"})
}

type DeleteMcpRequest struct {
	ID string `json:"id"`
}

// saveFlow 函数：保存 JSON 数据到 SQLite
func (s *Server) deleteMcp(id string) error {
	logrus.Debug("删除数据：", id)

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM mcps WHERE id = ?", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check flow existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("flow with id %s not found", id)
	}

	deleteSQL := `DELETE FROM mcps WHERE id=?`
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

func (s *Server) handleDeleteMcp(c *gin.Context) {
	data, err := c.GetRawData()
	logrus.Debug("delete data:", string(data))
	if err != nil {
		Error(c, http.StatusBadRequest, "数据请求错误")
		return
	}

	var body DeleteMcpRequest
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	logrus.Debug("delete id:", body.ID)
	err = s.deleteMcp(body.ID)
	if err != nil {
		logrus.Debug("删除失败：" + err.Error())
		Error(c, 100, "删除失败："+err.Error())
		return
	}

	logrus.Debug("删除数据成功")
	Success(c, "")
}

type ToolResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type CapsResponse struct {
	Prompts   []string       `json:"prompts"`
	Tools     []ToolResponse `json:"tools"`
	Resources []string       `json:"resources"`
}

func (s *Server) handleMcpCaps(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, "请求数据错误。")
		return
	}
	logrus.Debug(string(body))

	var item types.McpInfo
	err = json.Unmarshal(body, &item)
	if err != nil {
		logrus.Error("MCP数据解析错误:", err.Error())
		Error(c, http.StatusBadRequest, "MCP数据解析错误")
		return
	}

	mcp, err := components.CreateMcpSSEServer(&item)
	if err != nil {
		return
	}

	caps := CapsResponse{
		Prompts:   make([]string, 0),
		Tools:     make([]ToolResponse, 0),
		Resources: make([]string, 0),
	}

	ctx := context.Background()
	prompts, err := mcp.ListPrompts(ctx)
	if err != nil {
		logrus.Debug("Prompts:", err.Error())
		//return
	} else {
		logrus.Debug("Prompts:", prompts)
	}

	tools, err := mcp.ListTools(ctx, nil)
	if err != nil {
		logrus.Debug("Tools:", err.Error())
		//return
	} else {
		for _, t := range tools {
			t := ToolResponse{
				Name:        t.Name,
				Description: t.Desc,
			}
			caps.Tools = append(caps.Tools, t)
		}
	}

	resources, err := mcp.ListResources(ctx)
	if err != nil {
		logrus.Debug("Resources:", err.Error())
		//return
	} else {
		logrus.Debug("Resources:", resources)
	}
	mcp.Close()

	Success(c, caps)
}

// 根据 ID 获取模型记录
func (s *Server) getMcpByID(id string) (*types.McpInfo, error) {
	// 准备 SQL 查询语句
	query := "SELECT id, name, mcpType, url FROM mcps WHERE id = ?"
	row := s.db.QueryRow(query, id)

	// 解析查询结果
	var item types.McpInfo
	err := row.Scan(&item.ID,
		&item.Name,
		&item.McpType,
		&item.Url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no record found for ID %d", id)
		}
		return nil, err
	}

	return &item, nil
}

func (s *Server) CreateMcpServer(mcpId string) (*types.IMcpServer, error) {
	_, err := s.getMcpByID(mcpId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
