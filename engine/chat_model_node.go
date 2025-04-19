package engine

import (
	"context"
	"eino-script/parser"
	"fmt"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func (e *Engine) CreateChatModelNode(cfg *parser.NodeCfg) error {
	fmt.Println("CreateChatModelNode: ", *cfg)

	_ = e.g.AddChatModelNode(cfg.Name, &mockChatModel{}, compose.WithNodeName(cfg.Name))
	return nil
}

type mockChatModel struct{}

func (m *mockChatModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	return schema.AssistantMessage("the weather is good", nil), nil
}

func (m *mockChatModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	sr, sw := schema.Pipe[*schema.Message](0)
	go func() {
		defer sw.Close()
		sw.Send(schema.AssistantMessage("the weather is", nil), nil)
		sw.Send(schema.AssistantMessage("good", nil), nil)
	}()
	return sr, nil
}

func (m *mockChatModel) BindTools(tools []*schema.ToolInfo) error {
	panic("implement me")
}
