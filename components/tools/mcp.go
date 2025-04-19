package components

import (
	"context"
	"eino-script/types"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

type McpToolsConfig struct {
	Server    types.IMcpServer
	ToolsInfo []*schema.ToolInfo
}

type mcpToolHelper struct {
	server types.IMcpServer
	info   *schema.ToolInfo
}

func GetMcpTools(ctx context.Context, cfg *McpToolsConfig) ([]tool.BaseTool, error) {
	logrus.Debugf("GetMcpTools")
	tools := make([]tool.BaseTool, len(cfg.ToolsInfo))
	for i, toolInfo := range cfg.ToolsInfo {
		tools[i] = &mcpToolHelper{
			server: cfg.Server,
			info:   toolInfo,
		}
	}

	return tools, nil
}

func (t *mcpToolHelper) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return t.info, nil
}

func (t *mcpToolHelper) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	logrus.Debugf("Tools Streamable Run（%s）: %s", t.info.Name, argumentsInJSON)
	arg := make(map[string]any)
	err := json.Unmarshal([]byte(argumentsInJSON), &arg)
	if err != nil {
		return "", fmt.Errorf("unmarshal input fail: %w", err)
	}

	return t.server.CallTool(ctx, t.info.Name, argumentsInJSON)
}

func (t *mcpToolHelper) StreamableRun(ctx context.Context,
	argumentsInJSON string, opts ...tool.Option) (*schema.StreamReader[string], error) {
	logrus.Debugf("Tools Streamable Run: %s", argumentsInJSON)
	sr, sw := schema.Pipe[string](1)
	go func() {
		defer sw.Close()
		arg := make(map[string]any)
		err := json.Unmarshal([]byte(argumentsInJSON), &arg)
		if err != nil {
			sw.Send("", err)
			return
		}

		result, err := t.server.CallTool(ctx, t.info.Name, argumentsInJSON)
		if err != nil {
			sw.Send("", err)
			return
		}

		logrus.Debugf("Result: %s", result)
		sw.Send(result, nil)
	}()

	return sr, nil
}
