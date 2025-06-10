package server

import (
	"database/sql"
	"eino-script/engine"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

var ()

type Server struct {
	db          *sql.DB
	engineCache *LRUCache
}

func CreateServer() (*Server, error) {
	var err error
	s := &Server{
		engineCache: NewLRUCache(10),
	}

	err = s.initDB("flow.db")
	if err != nil {
		logrus.Error("数据库打开失败:", err)
		return nil, err
	}
	return s, nil
}

func (s *Server) close() {
	s.db.Close()
}

// 初始化 SQLite 数据库
func (s *Server) initDB(dbPath string) error {
	var err error
	s.db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// 创建表（如果不存在）
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS flows (
		id TEXT PRIMARY KEY,
		name TEXT,
		script TEXT NOT NULL
	);`
	_, err = s.db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

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

// saveFlow 函数：保存 JSON 数据到 SQLite
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

// 模拟大模型输出流
func generateMessages(c chan string, e *engine.Engine, msg string) {
	defer close(c)
	in := map[string]interface{}{
		"message": msg,
	}
	stream, err := e.Stream(in)
	if err != nil {
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
			return
		}
		logrus.Info(message)
		c <- fmt.Sprintf(`{"message":"%s"}`, message.Content)
	}
}

func StartServer() {
	s, err := CreateServer()
	if err != nil {
		logrus.Error(err)
		return
	}

	defer s.close()

	router := gin.Default()

	router.Use(corsMiddleware())

	//跨域操作
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "https://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/api/flow/list", s.handleGetFlowList)
	router.POST("/api/flow/get", s.handleGetFlow)
	router.POST("/api/flow/save", s.handleSaveFlow)
	router.POST("/api/flow/delete", s.handleDeleteFlow)
	router.POST("/api/flow/run", s.handlePlayFlow)
	router.POST("/api/flow/message", s.handleMessage)

	err = router.Run()
	if err != nil {
		return
	}
}

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24小时

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

type FlowListItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// 查询 ID 和 Name 列表
func fetchItems(db *sql.DB) ([]FlowListItem, error) {
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

func Error(c *gin.Context, code int, message string) {
	c.JSON(200, gin.H{
		"code": code,
		"data": gin.H{
			"message": message,
		},
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": data,
	})
}

func (s *Server) handleGetFlowList(c *gin.Context) {
	// 查询数据
	items, err := fetchItems(s.db)
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

	e, err := engine.CreateEngineByData(body, "json")
	if err != nil {
		logrus.Error(err)
		c.JSON(200, gin.H{
			"code": http.StatusInternalServerError,
			"data": gin.H{
				"message": err,
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

	messageChan := make(chan string)
	go generateMessages(messageChan, e, body.Message)

	for msg := range messageChan {
		fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
		flusher.Flush()
	}

	c.Writer.WriteHeaderNow()

}
