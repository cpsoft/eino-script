package engine

import (
	"eino-script/engine/types"
	"fmt"
	"github.com/sirupsen/logrus"
)

func (e *Engine) getNodeSourceId(id string) (string, error) {
	var node types.NodeInterface
	node, ok := e.nodes[id]
	if !ok {
		return "", fmt.Errorf("节点(%s)未找到。", id)
	}
	return node.GetSourceId()
}

func (e *Engine) getNodeTargetId(id string) (string, error) {
	var node types.NodeInterface
	node, ok := e.nodes[id]
	if !ok {
		return "", fmt.Errorf("节点(%s)未找到。", id)
	}
	return node.GetTargetId()
}

func (e *Engine) AddEdge(cfg *types.EdgeCfg) error {
	var sourceId string
	var targetId string
	var err error
	var sourceBranch types.BranchInterface
	var targetBranch types.BranchInterface

	//logrus.Infof("CreateEdge: %s -> %s", cfg.SourceNodeId, cfg.TargetNodeId)
	sourceNode, ok := e.nodes[cfg.SourceNodeId]
	if !ok {
		sourceBranch, ok = e.branchs[cfg.SourceNodeId]
		if !ok {
			return fmt.Errorf("没有找到（%s）对应的节点。", cfg.SourceNodeId)
		}
		sourceId = cfg.SourceNodeId
	} else {
		sourceId, err = sourceNode.GetSourceId()
		if err != nil {
			return fmt.Errorf("Edge参数错误。")
		}
	}

	targetNode, ok := e.nodes[cfg.TargetNodeId]
	if !ok {
		targetBranch, ok = e.branchs[cfg.TargetNodeId]
		if !ok {
			return fmt.Errorf("没有找到（%s）对应的节点。", cfg.SourceNodeId)
		}
		targetId = cfg.TargetNodeId
	} else {
		targetId, err = targetNode.GetTargetId()
		if err != nil {
			return fmt.Errorf("Edge参数错误。%s", err.Error())
		}
	}

	//如果源头是branch，需要链接源头的handle
	if sourceBranch != nil {
		err = sourceBranch.AddTargetNode(cfg.SourceHandle, targetId)
		if err != nil {
			return err
		}
		return nil
	}
	if targetBranch != nil {
		targetBranch.SetStartNode(cfg.SourceNodeId)
		return nil
	}
	logrus.Infof("CreateEdge: %s -> %s", sourceId, targetId)
	err = e.g.AddEdge(sourceId, targetId)
	if err != nil {
		return err
	}
	return nil
}
