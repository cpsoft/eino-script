package engine

import (
	"context"
	components "eino-script/components/template"
	"eino-script/parser"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/schema"
)

func (e *Engine) CreateMcpTemplateNode(cfg *parser.NodeCfg) error {
	fmt.Println("CreateMcpTemplateNode: ", *cfg)
	serverName, ok := cfg.Attrs["server"].(string)
	if serverName == "" {
		return errors.New("CreateMcpTemplateNode: name is required")
	}

	mcpServer, ok := e.mcps[serverName]
	if !ok {
		return errors.New("CreateMcpTemplateNode: mcp server not found")
	}

	messagesTemplate := make([]schema.MessagesTemplate, 0)
	messagesTemplate = append(messagesTemplate, schema.UserMessage("问题：{message}?"))
	promptCfg := &components.McpPromptConfig{
		Server: mcpServer,
		Name:   cfg.Name,
	}

	ctx := context.Background()
	pt := components.NewMcpPromptTemplate(
		ctx,
		promptCfg,
		schema.FString,
		messagesTemplate...,
	)

	_ = e.g.AddChatTemplateNode(cfg.Name, pt)
	return nil
}
