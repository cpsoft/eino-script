package components

import (
	"eino-script/engine/components/loaders"
	"eino-script/engine/types"
	"fmt"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

type LoaderNode struct {
	types.Node
	document.Loader
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

func CreateLoaderNode(cfg *types.NodeCfg, g *compose.Graph[any, any]) (types.NodeInterface, error) {
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

	node.Loader = loader

	err = g.AddLoaderNode(id, loader)
	if err != nil {
		return nil, err
	}

	return node, nil
}
