package server

import (
	"bytes"
	"eino-script/engine"
	"eino-script/engine/types"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"reflect"
)

// 定义与 JSON 对应的结构体
type SaveFlowRequestBody struct {
	ID   uint   `mapstructure:"id"` // 对应 JSON 中的 "id" 字段
	Name string `mapstructure:"name"`
}

// saveFlow 函数：保存 JSON 数据到 SQLite
func (s *Server) saveFlow(jsonData []byte) (uint, error) {
	// 解析 JSON 数据
	var req SaveFlowRequestBody
	v := viper.New()
	v.SetConfigType("json")
	if err := v.ReadConfig(bytes.NewReader(jsonData)); err != nil {
		return 0, err
	}

	if err := v.Unmarshal(&req); err != nil {
		return 0, err
	}

	if req.Name == "" {
		return 0, fmt.Errorf("工作流名字（%s）不能为空。", req.Name)
	}

	flow := types.FlowInfo{
		ID:     req.ID,
		Name:   req.Name,
		Script: string(jsonData),
	}

	// 插入或更新数据到数据库
	id, err := s.provider.SaveFlow(flow)
	if err != nil {
		return 0, fmt.Errorf("插入工作流失败. %s", err.Error())
	}

	return id, nil
}

func (s *Server) handleSaveFlow(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		logrus.Error(err)
		Error(c, http.StatusBadRequest, "请求数据错误："+string(body))
		return
	}

	logrus.Debug(string(body))
	id, err := s.saveFlow(body)
	if err != nil {
		logrus.Error(err)
		Error(c, http.StatusInternalServerError, "保存失败："+err.Error())
		return
	}

	Success(c, gin.H{
		"id": id,
	})
}

type DeleteFlowRequest struct {
	ID uint `json:"id"`
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
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Debug("delete id:", body.ID)
	s.provider.DeleteFlow(body.ID)
	logrus.Debug("删除数据成功")
	Success(c, "")
}

func (s *Server) handleGetFlowList(c *gin.Context) {
	// 查询数据
	flows, err := s.provider.GetFlowList()
	if err != nil {
		Error(c, 200, err.Error())
		return
	}
	Success(c, flows)
}

type GetFlowRequest struct {
	ID uint `json:"id"`
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
	flow, err := s.provider.GetFlow(body.ID)
	if err != nil {
		Error(c, 300, err.Error())
		return
	}
	Success(c, flow)
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
		"outmessage": msg,
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

		var text string
		switch msg := message.(type) {
		case string:
			//logrus.Debug("输出：" + msg)
			text = encodeToBase64(msg)
		case *schema.Message:
			//logrus.Debug("输出：" + msg.Content)
			text = encodeToBase64(msg.Content)
		case map[string]interface{}:
			//logrus.Debugf("map[string]interface: %v", msg)
			obj, _ := json.Marshal(msg)
			text = encodeToBase64(fmt.Sprintf("%+v", string(obj)))
			//c <- MessageOrError{Err: fmt.Errorf("未知输出类型")}
		default:
			v := reflect.TypeOf(message)
			logrus.Debug("未知输出类型：" + v.String())
			c <- MessageOrError{Err: fmt.Errorf("未知输出类型")}
			return
		}
		response := MessageOrError{
			Message: text,
		}
		c <- response
	}
}

func (s *Server) handlePlayFlow(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
	}

	id, err := s.saveFlow(body)
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
	s.engineCache.AddOrUpdate(e.Id(), e)

	Success(c, gin.H{
		"id": id,
	})
}

// 定义与 JSON 对应的结构体
type MessageRequestBody struct {
	ID      uint   `json:"id"`      // 对应 JSON 中的 "id" 字段
	Message string `json:"message"` // 对应 JSON 中的 "message" 字段
}

type MessageResponse struct {
	Message string `json:"message"`
}

func (s *Server) handleMessage(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		logrus.Error(err)
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	//logrus.Debug("Message:", string(data))

	// 解析 JSON 数据到结构体
	var body MessageRequestBody
	err = json.Unmarshal(data, &body)
	if err != nil {
		logrus.Debug("message 解析错误：", err.Error())
		Error(c, http.StatusBadRequest, err.Error())
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
