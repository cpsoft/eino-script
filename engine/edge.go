package engine

import (
	types2 "eino-script/engine/types"
	"fmt"
	"github.com/sirupsen/logrus"
)

func (e *Engine) getNodeSourceId(id string) (string, error) {
	var node types2.NodeInterface
	node, ok := e.nodes[id]
	if !ok {
		return "", fmt.Errorf("节点(%s)未找到。", id)
	}
	return node.GetSourceId()
}

func (e *Engine) getNodeTargetId(id string) (string, error) {
	var node types2.NodeInterface
	node, ok := e.nodes[id]
	if !ok {
		return "", fmt.Errorf("节点(%s)未找到。", id)
	}
	return node.GetTargetId()
}

func (e *Engine) AddEdge(cfg *types2.EdgeCfg) error {
	logrus.Infof("CreateEdge: %s -> %s", cfg.SourceNodeId, cfg.TargetNodeId)
	source, err := e.getNodeSourceId(cfg.SourceNodeId)
	if err != nil {
		return err
	}
	target, err := e.getNodeTargetId(cfg.TargetNodeId)
	if err != nil {
		return err
	}

	err = e.g.AddEdge(source, target)
	if err != nil {
		return err
	}
	return nil
}
