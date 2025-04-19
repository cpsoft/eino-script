package types

import (
	"context"
	"github.com/cloudwego/eino/schema"
)

type IMcpServer interface {
	Close()
	ListPrompts(ctx context.Context) (interface{}, error)
	ListTools(ctx context.Context, toolNameList []string) ([]*schema.ToolInfo, error)
	ListResources(ctx context.Context) (interface{}, error)
	CallTool(ctx context.Context, toolName string, argJson string) (string, error)
}
