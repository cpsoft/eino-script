package provider

import (
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DataProvider struct {
	db *gorm.DB
}

// 工作流
type Flow struct {
	gorm.Model
	Name   string
	Script string
}

// 模型
type Model struct {
	gorm.Model
	Name             string
	ModelType        string
	ModelName        string
	ApiKey           string
	ApiUrl           string
	MaxContextLength int
	StreamingEnabled bool
}

type Mcp struct {
	gorm.Model
	Name      string
	McpType   string
	Url       string
	Prompts   string
	Tools     string
	Resources string
}

type SessionMessage struct {
	gorm.Model
	SessionId string
	Name      string
	FlowId    uint
	Role      schema.RoleType
	Content   string
}

func NewDataProvider() (*DataProvider, error) {
	var err error
	provider := &DataProvider{}
	provider.db, err = gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("创建数据库失败：", err.Error())
	}

	err = provider.db.AutoMigrate(&Flow{})
	if err != nil {
		return nil, fmt.Errorf("创建工作流表失败：", err.Error())
	}

	err = provider.db.AutoMigrate(&Model{})
	if err != nil {
		return nil, fmt.Errorf("创建模型表失败：", err.Error())
	}

	err = provider.db.AutoMigrate(&Mcp{})
	if err != nil {
		return nil, fmt.Errorf("创建MCP表失败：", err.Error())
	}

	err = provider.db.AutoMigrate(&SessionMessage{})
	if err != nil {
		return nil, fmt.Errorf("创建对话Session表失败", err.Error())
	}

	return provider, nil
}

func (p *DataProvider) Close() {
}
