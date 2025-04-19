package engine

import (
	"context"
	"eino-script/components"
	"eino-script/parser"
	"eino-script/types"
	"fmt"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type Node struct {
	Name string
}

type Engine struct {
	ctx    context.Context
	g      *compose.Graph[map[string]any, *schema.Message]
	r      compose.Runnable[map[string]any, *schema.Message]
	s      *schema.StreamReader[*schema.Message]
	mcps   map[string]types.IMcpServer
	models map[string]model.ChatModel
	tools  map[string][]*schema.ToolInfo
	//nodes  []Node
}

func CreateEngine(cfg *parser.Config) (*Engine, error) {
	var err error
	e := &Engine{}
	e.ctx = context.Background()
	e.mcps = make(map[string]types.IMcpServer)
	e.tools = make(map[string][]*schema.ToolInfo)
	e.models = make(map[string]model.ChatModel)
	e.g = compose.NewGraph[map[string]any, *schema.Message]()

	for _, mcpCfg := range cfg.McpServers {
		fmt.Println("mcpCfg: ", mcpCfg.Name)
		switch mcpCfg.Type {
		case "SSEServer":
			server, err := components.CreateMcpSSEServer(&mcpCfg)
			if err != nil {
				return nil, err
			}
			e.mcps[mcpCfg.Name] = server
		}
	}

	for _, nodeCfg := range cfg.Nodes {
		fmt.Println("CreateEngine: ", nodeCfg.Name)
		switch nodeCfg.Type {
		case "ChatTemplate":
			err := e.CreateChatTemplateNode(&nodeCfg)
			if err != nil {
				return nil, err
			}
		case "McpTemplate":
			err := e.CreateMcpTemplateNode(&nodeCfg)
			if err != nil {
				return nil, err
			}
		case "ChatModel":
			err := e.CreateChatModelNode(&nodeCfg)
			if err != nil {
				return nil, err
			}
		case "McpToolNode":
			err := e.CreateMcpToolNode(&nodeCfg)
			if err != nil {
				return nil, err
			}
		case "OllamaChatModel":
			err := e.CreateOllamaChatModelNode(&nodeCfg)
			if err != nil {
				return nil, err
			}
		case "QwenChatModel":
			err := e.CreateQwenChatModelNode(&nodeCfg)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, edgeCfg := range cfg.Edges {
		fmt.Printf("CreateEdge: %s -> %s\n", edgeCfg.Src, edgeCfg.Dst)
		err = e.BindTools(edgeCfg.Src, edgeCfg.Dst)
		if err != nil {
			return nil, err
		}
		err = e.g.AddEdge(edgeCfg.Src, edgeCfg.Dst)
		if err != nil {
			return nil, err
		}
	}

	e.r, err = e.g.Compile(e.ctx, compose.WithMaxRunSteps(10))
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Engine) BindTools(modelName string, toolsName string) error {
	m, ok := e.models[modelName]
	if !ok {
		return nil
	}
	toolsNode, ok := e.tools[toolsName]
	if !ok {
		return nil
	}
	fmt.Printf("BindTools: %s <- %s\n", modelName, toolsName)
	err := m.BindTools(toolsNode)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) Invoke(in map[string]any) error {
	fmt.Println("Invoke: ", in)
	ret, err := e.r.Invoke(e.ctx, in)
	if err != nil {
		return err
	}
	fmt.Println(ret)
	return nil
}

func (e *Engine) Stream(in map[string]any) error {
	var err error
	fmt.Println("Stream: ", in)
	e.s, err = e.r.Stream(e.ctx, in)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) Close() {
	if e.s != nil {
		e.s.Close()
		e.s = nil
	}
	for _, mcp := range e.mcps {
		mcp.Close()
	}
	e.mcps = nil
}

func (e *Engine) Recv() (*schema.Message, error) {
	if e.s != nil {
		return e.s.Recv()
	}
	return nil, fmt.Errorf("Stream closed")
}
