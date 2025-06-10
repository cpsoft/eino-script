package types

type Node struct {
	NodeId   string
	NodeType NodeType
	Position Point
	Data     interface{}
}

type NodeInterface interface {
	Id() string
	Type() (NodeType, error)
	GetSourceId() (string, error)
	GetTargetId() (string, error)
}
