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
	ID               string `json:"id"`
	ModelName        string `json:"modelName"`
	ModelType        string `json:"modelType"`
	ApiKey           string `json:"apiKey"`
	ApiUrl           string `json:"apiUrl"`
	MaxContextLength int    `json:"maxContextLength"`
	StreamingEnabled bool   `json:"streamingEnabled"`
}
