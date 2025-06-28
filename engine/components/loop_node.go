package components

import (
	"eino-script/engine/types"
	"fmt"
	"github.com/cloudwego/eino/compose"
)

type LoopNode struct {
	types.Node
	sourceId   string
	conditions map[string]BranchCondition
}

func CreateLoop(cfg *types.NodeCfg, g *compose.Graph[any, any]) (types.BranchInterface, error) {
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
