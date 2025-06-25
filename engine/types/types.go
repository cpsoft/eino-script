package types

import "github.com/cloudwego/eino/schema"

type Node struct {
	NodeId   string
	NodeType NodeType
	Position Point
	Data     interface{}
}

type NodeInterface interface {
	Id() string
	Type() (NodeType, error)
	GetSourceId() (string, error)
	GetTargetId() (string, error)
}

type BranchInterface interface {
	Id() string
	SetStartNode(id string)
	AddTargetNode(handle string, id string) error
}

type SessionInterface interface {
	Id() uint
	GetMessages() []*schema.Message
	Close()
}

type ModelInfo struct {
	ID               uint   `json:"id"`        // 模型ID
	Name             string `json:"name"`      // 模型名称
	ModelType        string `json:"modelType"` // 模型类型 ollama、openai
	ModelName        string `json:"modelName"` // 大模型名称
	ApiKey           string `json:"apiKey"`
	ApiUrl           string `json:"apiUrl"`
	MaxContextLength int    `json:"maxContextLength"` // 最大token数
	StreamingEnabled bool   `json:"streamingEnabled"`
}

type McpToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type McpCaps struct {
	Prompts   []string      `json:"Prompts"`
	Tools     []McpToolInfo `json:"Tools"`
	Resources []string      `json:"Resources"`
}

type McpInfo struct {
	ID      uint   `json:"ID"`      // 模型ID
	Name    string `json:"Name"`    // 模型名称
	McpType string `json:"McpType"` // 模型类型 ollama、openai
	Url     string `json:"Url"`
	McpCaps
}

type FlowInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Script string `json:"script"`
}

type Callbacks interface {
	Callback_GetModelInfo(modelId uint) (*ModelInfo, error)
	Callback_CreateMcpServer(mcpId uint) (IMcpServer, error)
}
