package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SessionListRequest struct {
	FlowId uint `json:"flowId"`
}

func (s *Server) handleGetSessionList(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debug("handle get session list ", string(body))

	var info SessionListRequest
	err = json.Unmarshal(body, &info)
	if err != nil {
		logrus.Error("Unmarshal json failed", err)
		Error(c, http.StatusBadRequest, "历史记录请求ID格式错误。")
		return
	}

	items, err := s.provider.GetSessionList(info.FlowId)
	if err != nil {
		logrus.Error("GetSessionList failed", err)
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Debug("GetSessionList", items)
	Success(c, items)
}

type UpdateSessionRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) handleUpdateSession(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debug(string(body))

	var info UpdateSessionRequest
	err = json.Unmarshal(body, &info)
	if err != nil {
		logrus.Error("UpdateSession错误:", err.Error())
		Error(c, http.StatusBadRequest, "UpdateSession数据解析错误")
		return
	}

	// 插入或更新数据到数据库
	err = s.provider.UpdateSession(info.Id, info.Name)
	if err != nil {
		logrus.Debug("Session数据更新错误：", err.Error())
		Error(c, 200, "Session数据更新错误："+err.Error())
		return
	}

	Success(c, gin.H{})
}

type DeleteSessionRequest struct {
	ID string `json:"id"`
}

func (s *Server) handleDeleteSession(c *gin.Context) {
	data, err := c.GetRawData()
	logrus.Debug("delete data:", string(data))
	if err != nil {
		Error(c, http.StatusBadRequest, "数据请求错误")
		return
	}

	var body DeleteSessionRequest
	err = json.Unmarshal(data, &body)
	if err != nil {
		logrus.Error("Unmarshal json failed", err)
		Error(c, http.StatusBadRequest, "id解析错误。"+err.Error())
		return
	}

	logrus.Debug("delete id:", body.ID)
	err = s.provider.DeleteSession(body.ID)
	if err != nil {
		logrus.Debug("删除数据失败：", err.Error())
		Error(c, http.StatusInternalServerError, "删除数据失败："+err.Error())
		return
	}
	logrus.Debug("删除数据成功")
	Success(c, "")
}

type GetSessionRequest struct {
	ID string `json:"sessionId"`
}

func (s *Server) handleGetSessionMessages(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		Error(c, http.StatusBadRequest, "请求数据错误。")
		return
	}

	logrus.Debug("GetSession:", string(body))

	var req GetSessionRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		logrus.Error("Session数据解析错误:", err.Error())
		Error(c, http.StatusBadRequest, "Session数据解析错误")
		return
	}

	items, err := s.provider.GetSessionMessages(req.ID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "没有获取到指定会话")
		return
	}

	logrus.Debug("messages info:", *items)
	Success(c, *items)
}
