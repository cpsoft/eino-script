package parser

import (
	"bytes"
	"eino-script/engine/entity"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type FlowgramParser struct {
	v    *viper.Viper
	data map[string]any
}

func NewFlowgramParser(data []byte) *FlowgramParser {
	p := &FlowgramParser{}
	v := viper.New()
	v.SetConfigType("json")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil
	}

	p.data = v.AllSettings()
	return p
}

func parseStringVeriable(raw map[string]interface{}) (entity.Variable, error) {
	veriable := entity.VarString{}
	veriable.Type = "string"
	veriable.Default, _ = raw["default"].(string)

	return veriable.Variable, nil
}

func parseBoolVeriable(raw map[string]interface{}) (entity.Variable, error) {
	veriable := entity.VarBoolean{}
	veriable.Type = "boolean"
	veriable.Default, _ = raw["default"].(bool)
	return veriable.Variable, nil
}

func parseRefVeriable(raw map[string]interface{}) (entity.Variable, error) {
	var ok bool
	veriable := entity.VarRef{}
	veriable.Type = "ref"
	veriable.Content = make([]string, 0)
	content, ok := raw["content"].([]interface{})
	if !ok {
		return veriable.Variable, fmt.Errorf("Ref 类型解析错误")
	}
	for _, item := range content {
		veriable.Content = append(veriable.Content, item.(string))
	}

	return veriable.Variable, nil
}

func parseArrayVeriable(raw map[string]interface{}) (entity.Variable, error) {
	var err error
	veriable := entity.VarArray{}
	veriable.Type = "array"
	items, ok := raw["items"].(map[string]interface{})
	if !ok {
		return veriable.Variable, fmt.Errorf("Array items设置错误")
	}
	veriable.Items, err = parseVeriable(items)
	if err != nil {
		return veriable.Variable, err
	}
	return veriable.Variable, nil
}

func parseObjectVeriable(raw map[string]interface{}) (entity.Variable, error) {
	var err error
	veriable := entity.VarObject{}
	veriable.Type = "object"
	veriable.Properties = make(map[string]entity.Variable, 0)
	properties, ok := raw["properties"].(map[string]interface{})
	if ok {
		for k, v := range properties {
			propertie, ok := v.(map[string]interface{})
			if ok {
				veriable.Properties[k], err = parseVeriable(propertie)
				if err != nil {
					return veriable.Variable, err
				}
			}
		}
	}

	return veriable.Variable, nil
}

func parseVeriable(raw map[string]interface{}) (entity.Variable, error) {
	var veriable entity.Variable
	var err error

	Type, ok := raw["type"].(string)
	if !ok {
		return veriable, fmt.Errorf("变量解析错误")
	}

	switch Type {
	case "string":
		veriable, err = parseStringVeriable(raw)
	case "number":
		veriable = entity.Variable{Type: "number"}
	case "boolean":
		veriable, err = parseBoolVeriable(raw)
	case "object":
		veriable, err = parseObjectVeriable(raw)
	case "array":
		veriable, err = parseArrayVeriable(raw)
	case "ref":
		veriable, err = parseRefVeriable(raw)
	}

	if err != nil {
		return veriable, err
	}

	return veriable, nil
}

func parseConditionValue(raw map[string]interface{}) (entity.ConditionValue, error) {
	var err error
	value := entity.ConditionValue{}
	left, ok := raw["left"].(map[string]interface{})
	if !ok {
		return value, fmt.Errorf("Condition left value 格式错误")
	}
	value.Left, err = parseVeriable(left)
	if err != nil {
		return value, err
	}
	value.Operator, ok = raw["operator"].(string)
	if !ok {
		return value, fmt.Errorf("Condition operator 格式错误")
	}

	right, ok := raw["right"].(map[string]interface{})
	if ok {
		value.Right, err = parseVeriable(right)
		if err != nil {
			return value, err
		}
	}

	return value, nil
}

func parseConditions(list []interface{}) ([]entity.Condition, error) {
	var err error
	conditions := make([]entity.Condition, 0)
	for _, raw := range list {
		condition := entity.Condition{}
		c, ok := raw.(map[string]interface{})
		if !ok {
			return conditions, fmt.Errorf("Condition格式错误")
		}
		condition.Key, ok = c["key"].(string)
		if !ok {
			return conditions, fmt.Errorf("Condition Key设置错误")
		}
		value, ok := c["value"].(map[string]interface{})
		if !ok {
			return conditions, fmt.Errorf("Condition Value设置错误")
		}
		condition.Value, err = parseConditionValue(value)
		if err != nil {
			return conditions, err
		}
	}
	return conditions, nil
}

func parseNode(raw map[string]interface{}) (entity.Node, error) {
	var ok bool
	var err error
	node := entity.Node{}
	node.Id, ok = raw["id"].(string)
	if !ok {
		return node, fmt.Errorf("Node未设置id")
	}
	node.Type, ok = raw["type"].(string)
	if !ok {
		return node, fmt.Errorf("Node类型未设置")
	}

	data, ok := raw["data"].(map[string]any)
	if !ok {
		return node, fmt.Errorf("Node的data未解析")
	}

	for k, v := range data {
		switch k {
		case "title":
			node.Title = v.(string)
		case "description":
			node.Description = v.(string)
		case "outputs":
			node.Outputs, err = parseVeriable(v.(map[string]interface{}))
			if err != nil {
				return node, err
			}
		case "conditions":
			node.Conditions, err = parseConditions(v.([]interface{}))
			if err != nil {
				return node, err
			}
		default:

			break
		}
	}

	return node, nil
}

func parseNodes(raw map[string]interface{}) ([]entity.Node, error) {
	nodes := make([]entity.Node, 0)
	list, ok := raw["nodes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no nodes")
	}

	logrus.Debug("list %+v", list)

	for _, item := range list {
		n, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Node item error")
		}
		node, err := parseNode(n)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func parseEdge(data map[string]interface{}) (entity.Edge, error) {
	edge := entity.Edge{}
	return edge, nil
}

func parseEdges(raw map[string]interface{}) ([]entity.Edge, error) {
	edges := make([]entity.Edge, 0)
	list, ok := raw["edges"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no edges")
	}
	for _, item := range list {
		e, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Edge item error")
		}
		edge, err := parseEdge(e)
		if err != nil {
			return nil, err
		}
		edges = append(edges, edge)
	}
	return edges, nil
}

func (p *FlowgramParser) parseFlow() (flow entity.Flow, err error) {
	flow = entity.Flow{}
	flow.Id, _ = p.data["id"].(string)
	flow.Name, _ = p.data["name"].(string)

	flow.Nodes, err = parseNodes(p.data)
	if err != nil {
		return
	}

	flow.Edges, err = parseEdges(p.data)
	if err != nil {
		return
	}

	return
}

func ParseFlowgram(data []byte) (*entity.Flow, error) {
	p := NewFlowgramParser(data)
	flow, err := p.parseFlow()
	if err != nil {
		return nil, err
	}
	return &flow, nil
}
