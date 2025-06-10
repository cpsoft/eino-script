package engine

import (
	"eino-script/types"
	"fmt"
)

func CreateGeneralNode(cfg *types.NodeCfg) (*types.Node, error) {
	node := &types.Node{
		NodeId:   cfg.Id,
		NodeType: cfg.Type,
		Position: cfg.Position,
	}
	return node, nil
}

type StartNode struct {
	types.Node
}

func CreateStartNode(cfg *types.NodeCfg) (types.NodeInterface, error) {
	n, err := CreateGeneralNode(cfg)
	if err != nil {
		return nil, err
	}
	node := &StartNode{
		Node: *n,
	}
	return node, nil
}

func (s StartNode) Id() string {
	return s.NodeId
}

func (s StartNode) Type() (types.NodeType, error) {
	return s.NodeType, nil
}

func (s *StartNode) GetSourceId() (string, error) {
	return string(s.NodeType), nil
}

func (s *StartNode) GetTargetId() (string, error) {
	return "", fmt.Errorf("开始节点没有输入句柄")
}

type EndNode struct {
	types.Node
}

func CreateEndNode(cfg *types.NodeCfg) (types.NodeInterface, error) {
	n, err := CreateGeneralNode(cfg)
	if err != nil {
		return nil, err
	}
	node := &EndNode{
		Node: *n,
	}
	return node, nil
}

func (e EndNode) Id() string {
	return e.NodeId
}

func (e EndNode) Type() (types.NodeType, error) {
	return e.NodeType, nil
}

func (e EndNode) GetSourceId() (string, error) {
	return "", fmt.Errorf("结束节点没有输出句柄")
}

func (e EndNode) GetTargetId() (string, error) {
	return string(e.NodeType), nil
}
