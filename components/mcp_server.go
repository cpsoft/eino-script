package components

import (
	"context"
	"eino-script/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/pkg"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/cloudwego/eino/schema"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
)

type McpServer struct {
	client *client.Client
	caps   protocol.ServerCapabilities
}

func (m *McpServer) Close() {
	m.client.Close()
}

func (m *McpServer) ListPrompts(ctx context.Context) (interface{}, error) {
	if m.caps.Prompts == nil {
		return nil, errors.New("ListPrompts: mcp server does not support prompts")
	}
	prompts, err := m.client.ListPrompts(ctx)
	if err != nil {
		return nil, err
	}
	return prompts, nil
}

func convMessage(tool *protocol.Tool) (*schema.Message, error) {
	ret := &schema.Message{
		Role: schema.System,
	}
	msg, err := json.Marshal(tool)
	if err != nil {
		return nil, err
	}
	ret.Content = string(msg)
	return ret, nil
}

func (m *McpServer) ListTools(ctx context.Context, toolNameList []string) ([]*schema.ToolInfo, error) {
	var err error
	if m.caps.Tools == nil {
		return nil, errors.New("ListTools: mcp server does not support tools")
	}
	nameSet := make(map[string]struct{})
	for _, name := range toolNameList {
		nameSet[name] = struct{}{}
	}

	tools, err := m.client.ListTools(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*schema.ToolInfo, 0)
	for _, tool := range tools.Tools {
		if len(toolNameList) > 0 {
			if _, ok := nameSet[tool.Name]; !ok {
				continue
			}
		}

		marshaledInputSchema, err := json.Marshal(tool.InputSchema)
		if err != nil {
			return nil, err
		}
		inputSchema := &openapi3.Schema{}
		err = json.Unmarshal(marshaledInputSchema, &inputSchema)
		if err != nil {
			return nil, err
		}

		ret = append(ret, &schema.ToolInfo{
			Name:        tool.Name,
			Desc:        tool.Description,
			ParamsOneOf: schema.NewParamsOneOfByOpenAPIV3(inputSchema),
		})

	}
	return ret, nil
}

func (m *McpServer) CallTool(ctx context.Context, toolName string, argJson string) (string, error) {
	arg := make(map[string]any)
	err := json.Unmarshal([]byte(argJson), &arg)
	if err != nil {
		return "", err
	}
	result, err := m.client.CallTool(ctx, &protocol.CallToolRequest{
		Name:      toolName,
		Arguments: arg,
	})
	if err != nil {
		return "", err
	}

	logrus.Debugf("CallTool result: %v", result.Content)
	marshaledResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	if result.IsError {
		return "", fmt.Errorf("call mcp tool fail: %s", marshaledResult)
	}
	return string(marshaledResult), nil
}

func (m *McpServer) ListResources(ctx context.Context) (interface{}, error) {
	if m.caps.Resources == nil {
		return nil, errors.New("ListResources: mcp server does not support resources")
	}
	resources, err := m.client.ListResources(ctx)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

//func CreateMcpSSEServer(cfg *types.McpServerCfg) (types.IMcpServer, error) {
//	logrus.Debugf("CreateMcpSSEServer")
//	if cfg.Name == "" {
//		return nil, fmt.Errorf("Mcp server name is required")
//	}
//
//	url, ok := cfg.Attrs["server_url"].(string)
//	if !ok {
//		return nil, errors.New("server_url not found in attributes")
//	}
//
//	transportClient, err := transport.NewSSEClientTransport(url)
//	if err != nil {
//		return nil, err
//	}
//
//	client.WithLogger(pkg.DefaultLogger)
//
//	mcpClient, err := client.NewClient(transportClient, client.WithClientInfo(
//		protocol.Implementation{
//			Name:    "mcp client",
//			Version: protocol.Version,
//		}))
//	if err != nil {
//		return nil, err
//	}
//
//	caps := mcpClient.GetServerCapabilities()
//	fmt.Println(caps)
//	return &McpServer{mcpClient, caps}, nil
//}

func CreateMcpSSEServer(info *types.McpInfo) (types.IMcpServer, error) {
	logrus.Debugf("CreateMcpSSEServer")
	if info.Name == "" {
		return nil, fmt.Errorf("Mcp服务名字不能为空。")
	}
	transportClient, err := transport.NewSSEClientTransport(info.Url)
	if err != nil {
		return nil, err
	}

	client.WithLogger(pkg.DefaultLogger)

	mcpClient, err := client.NewClient(transportClient, client.WithClientInfo(
		protocol.Implementation{
			Name:    "mcp client",
			Version: protocol.Version,
		}))
	if err != nil {
		return nil, err
	}

	caps := mcpClient.GetServerCapabilities()
	fmt.Println(caps)
	return &McpServer{mcpClient, caps}, nil
}
