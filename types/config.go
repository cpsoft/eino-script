package types

// 定义枚举类型
type NodeType string

// 枚举值
const (
	NodeTypeStart        NodeType = "start"
	NodeTypeEnd          NodeType = "end"
	NodeTypeChatModel    NodeType = "chatmodel"
	NodeTypeChatTemplate NodeType = "chatTemplate"
	NodeTypeMcpTemplate  NodeType = "chatMcpTemplate"
	NodeTypeEvent        NodeType = "event"
)

type Point struct {
	X int
	Y int
}

type NodeCfg struct {
	Id       string                 `mapstructure:"id"`
	Type     NodeType               `mapstructure:"type"`
	Position Point                  `mapstructure:"position"`
	Name     string                 `mapstructure:"name"`
	Attrs    map[string]interface{} `mapstructure:",remain"`
}

type EdgeCfg struct {
	SourceNodeId string `mapstructure:"source"`
	TargetNodeId string `mapstructure:"target"`
	SourcePortId string `mapstructure:"sourcePortID"`
}

type McpServerCfg struct {
	Type  string                 `mapstructure:"type"`
	Name  string                 `mapstructure:"name"`
	Attrs map[string]interface{} `mapstructure:",remain"`
}

//type Config struct {
//	Tools         []NodeCfg      `mapstructure:"Tool"`
//	ChatModels    []NodeCfg      `mapstructure:"ChatModel"`
//	ChatTemplates []NodeCfg      `mapstructure:"ChatTemplate"`
//	Edges         []EdgeCfg      `mapstructure:"Edge"`
//	McpServers    []McpServerCfg `mapstructure:"McpServer"`
//}

type Config struct {
	Id    string    `mapstructure:"id"`
	Nodes []NodeCfg `mapstructure:"nodes"`
	Edges []EdgeCfg `mapstructure:"edges"`
}
