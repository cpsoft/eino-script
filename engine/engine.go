package engine

import (
	"bytes"
	"context"
	"eino-script/components"
	"eino-script/parser"
	"eino-script/types"
	"fmt"
	"github.com/cloudwego/eino-ext/callbacks/apmplus"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
)

type Node struct {
	Name string
}

type System struct {
	shutdown func(ctx context.Context) error
}

func InitSystem() (*System, error) {
	system := &System{}

	file, err := os.Open("config.toml")
	if err != nil {
		logrus.Warning("open file error:", err)
		return system, nil
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logrus.Warning("read file error, %s", err)
		return system, nil
	}

	v := viper.New()
	v.SetConfigType("toml")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	cfg := apmplus.Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		logrus.Warning("unmarshal config error, %s", err)
		return system, nil
	}

	cbh, shutdown, err := apmplus.NewApmplusHandler(&cfg)
	if err != nil {
		logrus.Warning("create apmplus handler error, %s", err)
		return system, nil
	}

	callbacks.AppendGlobalHandlers(cbh)
	system.shutdown = shutdown
	return system, nil
}

func (s System) Close() error {
	if s.shutdown != nil {
		s.shutdown(context.Background())
	}
	return nil
}

type Engine struct {
	ctx    context.Context
	g      *compose.Graph[map[string]any, *schema.Message]
	r      compose.Runnable[map[string]any, *schema.Message]
	s      *schema.StreamReader[*schema.Message]
	mcps   map[string]types.IMcpServer
	models map[string]model.ToolCallingChatModel
	tools  map[string][]*schema.ToolInfo
}

func CreateEngineByFile(filename string) (*Engine, error) {
	cfg, err := parser.ParserFile(filename)
	if err != nil {
		return nil, err
	}
	return CreateEngine(cfg)
}

func CreateEngineByData(data []byte) (*Engine, error) {
	cfg, err := parser.Parser(data)
	if err != nil {
		return nil, err
	}
	return CreateEngine(cfg)
}

func CreateEngine(cfg *parser.Config) (*Engine, error) {
	var err error
	e := &Engine{}
	e.ctx = context.Background()
	e.mcps = make(map[string]types.IMcpServer)
	e.tools = make(map[string][]*schema.ToolInfo)
	e.models = make(map[string]model.ToolCallingChatModel)
	e.g = compose.NewGraph[map[string]any, *schema.Message]()

	for _, mcpCfg := range cfg.McpServers {
		logrus.Infof("mcpCfg: %s", mcpCfg.Name)
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
		logrus.Infof("CreateEngine: %s", nodeCfg.Name)
		switch nodeCfg.Type {
		case "ChatTemplate":
			err = e.CreateChatTemplateNode(&nodeCfg)
		case "McpTemplate":
			err = e.CreateMcpTemplateNode(&nodeCfg)
		case "ChatModel":
			err = e.CreateChatModelNode(&nodeCfg)
		case "McpToolNode":
			err = e.CreateMcpToolNode(&nodeCfg)
		case "OllamaChatModel":
			err = e.CreateOllamaChatModelNode(&nodeCfg)
		case "QwenChatModel":
			err = e.CreateQwenChatModelNode(&nodeCfg)
		}
		if err != nil {
			return nil, err
		}
	}

	for _, edgeCfg := range cfg.Edges {
		err = e.AddEdge(&edgeCfg)
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
