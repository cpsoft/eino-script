package templates

import (
	"context"
	"eino-script/engine/types"
	"fmt"
	"github.com/cloudwego/eino-ext/components/prompt/mcp"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/mark3labs/mcp-go/client"
	extmcp "github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

func createStdioMcpClient(data map[string]interface{}) (*client.Client, error) {
	return nil, fmt.Errorf("暂不支持stdio模式。")
}

func createSSEMcpClient(data map[string]interface{}) (*client.Client, error) {
	baseUrl, ok := data["sseurl"].(string)
	if !ok {
		return nil, fmt.Errorf("SSE模式下，BaseURL没有配置")
	}
	cli, err := client.NewSSEMCPClient(baseUrl)
	if err != nil {
		return nil, err
	}
	err = cli.Start(context.Background())
	if err != nil {
		return nil, err
	}
	initRequest := extmcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = extmcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = extmcp.Implementation{
		Name:    "aiflow-client",
		Version: "1.0.0",
	}
	_, err = cli.Initialize(context.Background(), initRequest)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func CreateMcpTemplateNode(cfg *types.NodeCfg) (prompt.ChatTemplate, error) {
	var err error
	logrus.Infof("CreateMcpTemplateNode: %+v", *cfg)

	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in attrs")
	}

	logrus.Debug("attrs:", data)

	mode, ok := data["mode"].(string)
	if !ok {
		return nil, fmt.Errorf("MCP 模式没有配置。")
	}

	var cli *client.Client
	switch mode {
	case "stdio":
		cli, err = createStdioMcpClient(data)
		if err != nil {
			return nil, err
		}
		break
	case "sse":
		cli, err = createSSEMcpClient(data)
		if err != nil {
			return nil, err
		}
		break
	default:
		return nil, fmt.Errorf("未知MCP类型：%s。", mode)
	}

	tpl, err := mcp.NewPromptTemplate(context.Background(), &mcp.Config{
		Cli:  cli,
		Name: "mcpclient",
	})

	if err != nil {
		return nil, err
	}

	result, err := tpl.Format(context.Background(), map[string]interface{}{"persona": "上海到北京怎么走"})
	if err != nil {
		return nil, err
	}

	logrus.Debug("mcp result  :   ", result)
	return tpl, nil
}
