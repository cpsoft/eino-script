package types

type NodeCfg struct {
	Type  string                 `mapstructure:"type"`
	Name  string                 `mapstructure:"name"`
	Attrs map[string]interface{} `mapstructure:",remain"`
}

type EdgeCfg struct {
	Src string `mapstructure:"src"`
	Dst string `mapstructure:"dst"`
}

type McpServerCfg struct {
	Type  string                 `mapstructure:"type"`
	Name  string                 `mapstructure:"name"`
	Attrs map[string]interface{} `mapstructure:",remain"`
}

type Config struct {
	Tools         []NodeCfg      `mapstructure:"Tool"`
	ChatModels    []NodeCfg      `mapstructure:"ChatModel"`
	ChatTemplates []NodeCfg      `mapstructure:"ChatTemplate"`
	Edges         []EdgeCfg      `mapstructure:"Edge"`
	McpServers    []McpServerCfg `mapstructure:"McpServer"`
}
