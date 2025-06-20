package engine

import (
	"bytes"
	"context"
	"eino-script/engine/types"
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
	id        uint
	ctx       context.Context
	callbacks types.Callbacks
	g         *compose.Graph[any, any]
	r         compose.Runnable[any, any]
	s         *schema.StreamReader[*schema.Message]
	mcps      map[string]types.IMcpServer
	models    map[string]model.ToolCallingChatModel
	tools     map[string][]*schema.ToolInfo
	nodes     map[string]types.NodeInterface
}

func CreateEngineByFile(callbacks types.Callbacks, filename string) (*Engine, error) {
	cfg, err := ParserFile(filename)
	if err != nil {
		return nil, err
	}
	return CreateEngine(callbacks, cfg)
}

func CreateEngineByData(callbacks types.Callbacks, data []byte, format string) (*Engine, error) {
	cfg, err := Parser(data, format)
	if err != nil {
		return nil, err
	}
	return CreateEngine(callbacks, cfg)
}

func CreateEngine(cb types.Callbacks, cfg *types.Config) (*Engine, error) {
	var err error
	e := &Engine{}
	e.id = cfg.Id
	e.callbacks = cb
	e.ctx = context.Background()
	e.mcps = make(map[string]types.IMcpServer)
	e.tools = make(map[string][]*schema.ToolInfo)
	e.models = make(map[string]model.ToolCallingChatModel)
	e.nodes = make(map[string]types.NodeInterface)
	e.g = compose.NewGraph[any, any]()

	//handler := callbacks.NewHandlerBuilder().OnStartFn(
	//	func(ctx context.Context,
	//		info *callbacks.RunInfo,
	//		input callbacks.CallbackInput) context.Context {
	//		logrus.Debugf("%+v", *info)
	//		return ctx
	//	}).Build()
	//e.ctx = callbacks.InitCallbacks(context.Background(), nil, handler)

	err = e.CreateNodes(&cfg.Nodes)
	if err != nil {
		return nil, err
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

func (e *Engine) Id() uint {
	return e.id
}

func (e *Engine) Invoke(in map[string]any) (any, error) {
	return e.r.Invoke(e.ctx, in)
}

func (e *Engine) Stream(in map[string]any) (*schema.StreamReader[any], error) {
	return e.r.Stream(e.ctx, in)
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
