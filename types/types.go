package types

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

type ModelInfo struct {
	ID               string `json:"id"`        // 模型ID
	Name             string `json:"name"`      // 模型名称
	ModelType        string `json:"modelType"` // 模型类型 ollama、openai
	ModelName        string `json:"modelName"` // 大模型名称
	ApiKey           string `json:"apiKey"`
	ApiUrl           string `json:"apiUrl"`
	MaxContextLength int    `json:"maxContextLength"` // 最大token数
	StreamingEnabled bool   `json:"streamingEnabled"`
}
