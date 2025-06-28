package entity

type Variable struct {
	Type string
}

type VarString struct {
	Variable
	Default string
}

type VarBoolean struct {
	Variable
	Default bool
}

type VarRef struct {
	Variable
	Content []string
}

type VarArray struct {
	Variable
	Items Variable //数组对象的内容定义
}

type VarObject struct {
	Variable
	Properties map[string]Variable
}

type ConditionValue struct {
	Left     Variable
	Operator string
	Right    Variable
}

type Condition struct {
	Key   string
	Value ConditionValue
}

type Node struct {
	Id          string
	Type        string
	Title       string
	Description string
	Outputs     Variable
	Conditions  []Condition
	Custom      Variable
}

type Edge struct {
	Id           string
	Source       string
	Target       string
	SourceHandle string
	TargetHandle string
}

type Flow struct {
	Id    string
	Name  string
	Nodes []Node
	Edges []Edge
}
