package server

import (
	"context"
	"database/sql"
	"eino-script/components"
	"eino-script/engine/types"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) handleGetMcpList(c *gin.Context) {
	// 查询数据
	items, err := s.provider.GetMcpList()
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

	var info types.McpInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		logrus.Error("MCP数据解析错误:", err.Error())
		Error(c, http.StatusBadRequest, "MCP数据解析错误")
		return
	}

	caps, err := GetMcpCapsFromServer(&info)
	if err != nil {
		Error(c, 200, "获取Mcp能力失败。"+err.Error())
		return
	}
	info.McpCaps = *caps

	// 插入或更新数据到数据库
	id, err := s.provider.SaveMcp(&info)
	if err != nil {
		logrus.Debug("MCP插入错误：", err.Error())
		Error(c, 200, "插入MCP错误："+err.Error())
		return
	}

	Success(c, gin.H{"id": id})
}

type DeleteMcpRequest struct {
	ID uint `json:"id"`
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
	err = s.provider.DeleteMcp(body.ID)
	if err != nil {
		logrus.Debug("删除数据失败：", err.Error())
		Error(c, http.StatusInternalServerError, "删除数据失败："+err.Error())
		return
	}
	logrus.Debug("删除数据成功")
	Success(c, "")
}

func GetMcpCapsFromServer(info *types.McpInfo) (*types.McpCaps, error) {
	mcp, err := components.CreateMcpSSEServer(info)
	if err != nil {
		return nil, err
	}

	caps := types.McpCaps{
		Prompts:   make([]string, 0),
		Tools:     make([]types.McpToolInfo, 0),
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
			t := types.McpToolInfo{
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
	return &caps, nil
}

type GetMcpRequest struct {
	ID uint `json:"id"`
}

func (s *Server) handleGetMcp(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, "请求数据错误。")
		return
	}

	logrus.Debug("GetMcp:", string(body))

	var req GetMcpRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		logrus.Error("MCP数据解析错误:", err.Error())
		Error(c, http.StatusBadRequest, "MCP数据解析错误")
		return
	}

	mcp, err := s.provider.GetMcp(req.ID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "没有获取到指定MCP")
		return
	}

	info := types.McpInfo{
		ID:      mcp.ID,
		Name:    mcp.Name,
		McpType: mcp.McpType,
		Url:     mcp.Url,
	}

	err = json.Unmarshal([]byte(mcp.Prompts), &info.Prompts)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Prompts解析错误。")
		return
	}

	err = json.Unmarshal([]byte(mcp.Tools), &info.Tools)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Tools解析错误。")
		return
	}

	err = json.Unmarshal([]byte(mcp.Resources), &info.Resources)
	if err != nil {
		Error(c, http.StatusInternalServerError, "Resources解析错误。")
		return
	}

	logrus.Debug("mcp info:", info)
	Success(c, info)
}

func (s *Server) handleMcpCaps(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, "请求数据错误。")
		return
	}
	logrus.Debug(string(body))

	var info types.McpInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		logrus.Error("MCP数据解析错误:", err.Error())
		Error(c, http.StatusBadRequest, "MCP数据解析错误")
		return
	}

	caps, err := GetMcpCapsFromServer(&info)
	if err != nil {
		Error(c, 200, "获取Mcp能力失败。"+err.Error())
		return
	}

	Success(c, caps)
}

// 根据 ID 获取模型记录
func (s *Server) getMcpByID(id uint) (*types.McpInfo, error) {
	// 准备 SQL 查询语句
	mcp, err := s.provider.GetMcp(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no record found for ID %d", id)
		}
		return nil, err
	}

	info := types.McpInfo{
		ID:      mcp.ID,
		Name:    mcp.Name,
		McpType: mcp.McpType,
		Url:     mcp.Url,
	}

	return &info, nil
}

func (s *Server) Callback_CreateMcpServer(mcpId uint) (types.IMcpServer, error) {

	mcp, ok := s.mcpCache.Get(mcpId)
	if ok {
		return mcp, nil
	}

	info, err := s.getMcpByID(mcpId)
	if err != nil {
		return nil, err
	}

	mcp, err = components.CreateMcpSSEServer(info)
	if err != nil {
		return nil, err
	}

	s.mcpCache.AddOrUpdate(mcpId, mcp)

	return mcp, nil
}
