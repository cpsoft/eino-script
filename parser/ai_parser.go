package parser

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
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
	logrus.Debugf("flow script：%+v", cfg)
	return &cfg, nil
}

func ParserFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		logrus.Error("open file error:", err)
		return nil, fmt.Errorf("open file error, %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logrus.Errorf("read file error, %s", err)
		return nil, fmt.Errorf("read file error, %s", err)
	}
	return Parser(data)
}
