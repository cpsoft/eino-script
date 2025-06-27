package engine

import (
	"eino-script/engine/loaders"
	"eino-script/engine/types"
	"fmt"
)

type LoaderNode struct {
	types.Node
	loader types.LoaderInterface
}

func (l LoaderNode) Id() string {
	return l.NodeId
}

func (l LoaderNode) Type() (types.NodeType, error) {
	return l.NodeType, nil
}

func (l *LoaderNode) GetTargetId() (string, error) {
	return l.NodeId, nil
}

func (l *LoaderNode) GetSourceId() (string, error) {
	return l.NodeId, nil
}

func (e *Engine) CreateLoaderNode(cfg *types.NodeCfg) (types.NodeInterface, error) {
	n, err := CreateGeneralNode(cfg)
	if err != nil {
		return nil, err
	}

	node := &LoaderNode{
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

	loader, err := loaders.CreateLoaderNode(data)
	if err != nil {
		return nil, err
	}

	node.loader = loader

	einoLoader, err := node.loader.GetEinoNode()
	if err != nil {
		return nil, err
	}

	//Todo: 对齐问题还需要处理
	err = e.g.AddLoaderNode(id, einoLoader)
	if err != nil {
		return nil, err
	}

	return node, nil
}
