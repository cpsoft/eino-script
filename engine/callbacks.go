package engine

import "eino-script/types"

type Callbacks interface {
	GetModelInfo(modelName string) (*types.ModelInfo, error)
	CreateMcpServer(mcpId string) (*types.IMcpServer, error)
}
