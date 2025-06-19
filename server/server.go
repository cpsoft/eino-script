package server

import (
	"eino-script/server/provider"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	provider    *provider.DataProvider
	engineCache *LRUCache
}

func CreateServer() (*Server, error) {
	var err error
	s := &Server{
		engineCache: NewLRUCache(10),
	}

	s.provider, err = provider.NewDataProvider()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) close() {
	s.provider.Close()
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

	router.GET("/api/model/list", s.handleGetModelList)
	router.POST("/api/model/save", s.handleSaveModel)
	router.POST("api/model/delete", s.handleDeleteModel)
	router.POST("/api/model/chatmodel/list", s.handleChatModelList)

	router.GET("/api/mcp/list", s.handleGetMcpList)
	router.POST("/api/mcp/get", s.handleGetMcp)
	router.POST("/api/mcp/save", s.handleSaveMcp)
	router.POST("api/mcp/delete", s.handleDeleteMcp)
	router.POST("api/mcp/getCaps", s.handleMcpCaps)

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
