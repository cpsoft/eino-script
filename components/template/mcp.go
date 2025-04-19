package components

import (
	"context"
	"eino-script/types"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type McpPromptConfig struct {
	Server types.IMcpServer
	Name   string
}

type chatTemplate struct {
	name       string
	formatType schema.FormatType
	templates  []schema.MessagesTemplate
	server     types.IMcpServer
}

func (t *chatTemplate) Format(ctx context.Context,
	vs map[string]any, _ ...prompt.Option) (result []*schema.Message, err error) {
	//defer func() {
	//	if err != nil {
	//		_ = callbacks.OnError(ctx, err)
	//	}
	//}()
	//
	//ctx = callbacks.OnStart(ctx, &prompt.CallbackInput{
	//	Variables: vs,
	//	Templates: t.templates,
	//})
	//
	//result = make([]*schema.Message, 0)
	//for _, template := range t.templates {
	//	msgs, err := template.Format(ctx, vs, t.formatType)
	//	if err != nil {
	//		return nil, err
	//	}
	//	result = append(result, msgs...)
	//}
	//
	//tools, err := t.server.ListTools(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(tools)
	//result = append(result, tools...)
	//
	//_ = callbacks.OnEnd(ctx, &prompt.CallbackOutput{
	//	Result:    result,
	//	Templates: t.templates,
	//})
	return result, nil
}

func (t *chatTemplate) GetType() string {
	return "MCP"
}

func (t *chatTemplate) IsCallbacksEnabled() bool {
	return true
}

func NewMcpPromptTemplate(ctx context.Context,
	conf *McpPromptConfig,
	formatType schema.FormatType,
	templates ...schema.MessagesTemplate) prompt.ChatTemplate {
	return &chatTemplate{
		server:     conf.Server,
		name:       conf.Name,
		formatType: formatType,
		templates:  templates,
	}
}
