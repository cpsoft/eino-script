package components

import (
	"context"
	"eino-script/engine/types"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

type BranchCondition struct {
	Id     string
	NodeId string
	Type   types.ConditionType
	Name   string
	Script string
}

type BranchNode struct {
	types.Node
	sourceId   string
	conditions map[string]BranchCondition
}

func (bn *BranchNode) Id() string {
	return bn.NodeId
}

func (bn *BranchNode) SetStartNode(id string) {
	logrus.Debugf("Branch(%s) sourceId (%s)", bn.Id(), id)
	bn.sourceId = id
}

func (bn *BranchNode) AddTargetNode(handle string, id string) error {
	logrus.Debugf("Branch(%s) targetId[%s] (%s)", bn.Id(), handle, id)
	condition, ok := bn.conditions[handle]
	if !ok {
		return fmt.Errorf("条件（%s）不存在。", handle)
	}
	condition.NodeId = id
	bn.conditions[handle] = condition
	return nil
}

func (bn BranchNode) CommonBranch(ctx context.Context, input *schema.StreamReader[any]) (endNode string, err error) {
	logrus.Debug("CommonBranch")
	defer input.Close()
	for _, condition := range bn.conditions {
		logrus.Debug(condition.Type)
		switch condition.Type {
		case types.ConditionTypeNeedTools:
			data, err := input.Recv()
			if err != nil {
				return "", err
			}
			msg, ok := data.(*schema.Message)
			if ok && len(msg.ToolCalls) > 0 {
				logrus.Debugf("CommonBranch get toolCalls")
				return condition.NodeId, nil
			}
			break
		case types.ConditionTypeDefault:
			logrus.Debug("default:", condition.NodeId)
			return condition.NodeId, nil
		}
	}
	return "", nil
}

func CreateBranch(cfg *types.NodeCfg) (types.BranchInterface, error) {
	n, err := CreateGeneralNode(cfg)
	if err != nil {
		return nil, err
	}
	node := &BranchNode{
		Node:       *n,
		conditions: make(map[string]BranchCondition),
	}

	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Branch的data参数错误。")
	}

	list, ok := data["conditions"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Branch的conditions参数错误。")
	}

	for _, item := range list {
		info, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Branch的condition参数错误。")
		}
		condition := BranchCondition{}
		condition.Id, ok = info["id"].(string)
		if !ok {
			return nil, fmt.Errorf("Branch的conditions参数错误。")
		}
		condition.Name, ok = info["name"].(string)
		if !ok {
			return nil, fmt.Errorf("Branch的conditions参数错误。")
		}
		conditionType, ok := info["type"].(string)
		if !ok {
			return nil, fmt.Errorf("Branch的conditions参数错误。")
		}
		condition.Type = types.ConditionType(conditionType)
		condition.Script, ok = info["script"].(string)

		node.conditions[condition.Id] = condition
	}

	return node, nil
}

func BranchsInit(g *compose.Graph[any, any], branchs *map[string]types.BranchInterface) error {
	for _, info := range *branchs {
		b, ok := info.(*BranchNode)
		if !ok {
			return fmt.Errorf("Branch节点错误。")
		}
		targets := make(map[string]bool)
		for _, handle := range b.conditions {
			targets[handle.NodeId] = true
		}
		logrus.Debug("branch(%s):", b.Id(), b.sourceId)
		err := g.AddBranch(b.sourceId, compose.NewStreamGraphBranch(b.CommonBranch, targets))
		if err != nil {
			return err
		}
	}
	return nil
}
