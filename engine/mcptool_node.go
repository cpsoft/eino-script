package engine

import (
	"eino-script/engine/nodes"
	"eino-script/engine/types"
	"fmt"
	"github.com/cloudwego/eino/compose"
)

type McpToolNode struct {
	types.Node
}

func (e McpToolNode) Id() string {
	return e.NodeId
}

func (e McpToolNode) Type() (types.NodeType, error) {
	return e.NodeType, nil
}

func (cm *McpToolNode) GetTargetId() (string, error) {
	return cm.NodeId, nil
}

func (cm *McpToolNode) GetSourceId() (string, error) {
	return cm.NodeId, nil
}

func (e *Engine) CreateMcpToolNode(cfg *types.NodeCfg) (types.NodeInterface, error) {
	n, err := CreateGeneralNode(cfg)
	if err != nil {
		return nil, err
	}

	node := &McpToolNode{
		Node: *n,
	}

	id := node.Id()
	if id == "" {
		return nil, err
	}

	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in attrs")
	}

	mcpId, ok := data["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("McpTool的mcpId没有设置。")
	}

	server, err := e.callbacks.Callback_CreateMcpServer(uint(mcpId))
	if err != nil {
		return nil, err
	}

	mcpTool, err := nodes.CreateMcpToolNode(cfg, server)
	if err != nil {
		return nil, err
	}

	//Todo: 对齐问题还需要处理
	err = e.g.AddToolsNode(id, mcpTool, compose.WithOutputKey("outmessage"))
	if err != nil {
		return nil, err
	}

	return node, nil
}
