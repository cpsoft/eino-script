package parser

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
)

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
	Nodes      []NodeCfg      `mapstructure:"Node"`
	Edges      []EdgeCfg      `mapstructure:"Edge"`
	McpServers []McpServerCfg `mapstructure:"McpServer"`
}

func Parser(data []byte) (*Config, error) {
	var cfg Config
	v := viper.New()
	v.SetConfigType("toml")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	fmt.Println("数据：", cfg)
	return &cfg, nil
}
