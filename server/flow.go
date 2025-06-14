package server

import (
	"database/sql"
	"eino-script/engine"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// 定义与 JSON 对应的结构体
type SaveFlowRequestBody struct {
	ID   string `json:"id"` // 对应 JSON 中的 "id" 字段
	Name string `json:"name"`
}

// saveFlow 函数：保存 JSON 数据到 SQLite
func (s *Server) saveFlow(jsonData []byte) error {
	// 解析 JSON 数据
	var flow SaveFlowRequestBody
	err := json.Unmarshal(jsonData, &flow)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	logrus.Debug("保存数据")
	if flow.ID == "" || flow.Name == "" {
		return fmt.Errorf("工作流ID（%s）或名字（%s）不能为空。", flow.ID, flow.Name)
	}

	// 插入或更新数据到数据库
	upsertSQL := `
	INSERT INTO flows (id, name, script) 
	VALUES (?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET 
		name = excluded.name, 
		script = excluded.script
	`
	_, err = s.db.Exec(upsertSQL, flow.ID, flow.Name, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to insert or update data: %w", err)
	}

	return nil
}

// deleteFlow 函数：保存 JSON 数据到 SQLite
func (s *Server) deleteFlow(id string) error {
	logrus.Debug("删除数据：", id)

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM flows WHERE id = ?", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check flow existence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("flow with id %s not found", id)
	}

	deleteSQL := `DELETE FROM flows WHERE id=?`
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

type MessageOrError struct {
	Message string
	Err     error
}

func encodeToBase64(content string) string {
	// 将字符串转换为字节数组
	contentBytes := []byte(content)

	// 使用 Base64 编码
	base64Encoded := base64.StdEncoding.EncodeToString(contentBytes)

	return base64Encoded
}

// 模拟大模型输出流
func generateMessages(c chan MessageOrError, e *engine.Engine, msg string) {
	defer close(c)
	in := map[string]interface{}{
		"message": msg,
	}
	stream, err := e.Stream(in)
	if err != nil {
		c <- MessageOrError{Err: err}
		return
	}
	defer stream.Close()

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			logrus.Error("recv failed: %v", err)
			c <- MessageOrError{Err: err}
			return
		}
		text := encodeToBase64(message.Content)
		response := MessageOrError{
			Message: text,
		}
		c <- response
	}
}

type FlowListItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// 查询 ID 和 Name 列表
func fetchFlowItems(db *sql.DB) ([]FlowListItem, error) {
	query := `SELECT id, name FROM flows`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []FlowListItem
	for rows.Next() {
		var item FlowListItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	return items, nil
}

func (s *Server) handleGetFlowList(c *gin.Context) {
	// 查询数据
	items, err := fetchFlowItems(s.db)
	logrus.Debug("items:", items)
	if err != nil {
		Error(c, 300, "读取数据失败："+err.Error())
		return
	}

	Success(c, items)

}

// 根据 ID 查询单条记录
func fetchItemByID(db *sql.DB, id string) (string, error) {
	query := `SELECT script FROM flows WHERE id = ?`
	row := db.QueryRow(query, id)

	var script string
	if err := row.Scan(&script); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("item with id %s not found", id)
		}
		return "", fmt.Errorf("failed to scan row: %w", err)
	}

	return script, nil
}

type GetFlowRequest struct {
	ID string `json:"id"`
}

func (s *Server) handleGetFlow(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		Error(c, 300, err.Error())
		return
	}
	// 解析 JSON 数据到结构体
	var body GetFlowRequest
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	// 打印解析结果（或进行其他逻辑处理）
	logrus.Infof("Parsed Body: ID=%s, Message=%s\n", body.ID)
	script, err := fetchItemByID(s.db, body.ID)
	if err != nil {
		Error(c, 300, err.Error())
		return
	}
	Success(c, script)
}

func (s *Server) handleSaveFlow(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
	}

	logrus.Debug(string(body))
	err = s.saveFlow(body)
	if err != nil {
		logrus.Error(err)
		c.JSON(200, gin.H{
			"code": "100",
			"data": gin.H{
				"message": "保存失败:" + err.Error(),
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"message": "保存成功",
		},
	})
}

type DeleteFlowRequest struct {
	ID string `json:"id"`
}

func (s *Server) handleDeleteFlow(c *gin.Context) {
	data, err := c.GetRawData()
	logrus.Debug("delete data:", string(data))
	if err != nil {
		Error(c, http.StatusBadRequest, "数据请求错误")
		return
	}

	var body DeleteFlowRequest
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	logrus.Debug("delete id:", body.ID)
	err = s.deleteFlow(body.ID)
	if err != nil {
		logrus.Debug("删除失败：" + err.Error())
		Error(c, 100, "删除失败："+err.Error())
		return
	}

	logrus.Debug("删除数据成功")
	Success(c, "")
}

func (s *Server) handlePlayFlow(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
	}

	err = s.saveFlow(body)
	if err != nil {
		logrus.Error(err)
		Error(c, 100, "保存失败")
		return
	}

	e, err := engine.CreateEngineByData(s, body, "json")
	if err != nil {
		logrus.Error(err)
		c.JSON(200, gin.H{
			"code": http.StatusInternalServerError,
			"data": gin.H{
				"message": err.Error(),
			},
		})
		return
	}
	s.engineCache.AddOrUpdate(e)

	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"message": "保存成功",
		},
	})
}

// 定义与 JSON 对应的结构体
type MessageRequestBody struct {
	ID      string `json:"id"`      // 对应 JSON 中的 "id" 字段
	Message string `json:"message"` // 对应 JSON 中的 "message" 字段
}

type MessageResponse struct {
	Message string `json:"message"`
}

func (s *Server) handleMessage(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
		return
	}

	// 解析 JSON 数据到结构体
	var body MessageRequestBody
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	// 打印解析结果（或进行其他逻辑处理）
	logrus.Infof("Parsed Body: ID=%s, Message=%s\n", body.ID, body.Message)

	e, ok := s.engineCache.Get(body.ID)
	if !ok {
		logrus.Error("工作流不存在")
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "工作流不存在",
		})
		return
	}

	// 设置 SSE 所需的响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error": "Streaming unsupported",
			},
		)
		return
	}

	messages := make(chan MessageOrError)
	go generateMessages(messages, e, body.Message)

	for msg := range messages {
		// 在发生错误时，通过 SSE 格式发送错误消息
		if msg.Err != nil {
			_, _ = fmt.Fprintf(c.Writer, "error: %s\n\n", msg.Err.Error())
			flusher.Flush()
			c.Writer.CloseNotify()
			return
		}
		response := MessageResponse{
			Message: msg.Message,
		}
		//fmt.Printf("data: %s\n\n", msg)
		jsonData, err := json.Marshal(response)
		if err != nil {
			_, _ = fmt.Fprintf(c.Writer, "error: %s\n\n", msg.Err.Error())
			flusher.Flush()
			c.Writer.CloseNotify()
			return
		}
		fmt.Println(string(jsonData))
		_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", jsonData)
		flusher.Flush()
	}

	c.Writer.WriteHeaderNow()

}
