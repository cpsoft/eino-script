package parser

import (
	"bytes"
	"eino-script/engine/entity"
	engine "eino-script/engine/types"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
)

func Parser(data []byte, format string) (*engine.Config, error) {
	var cfg engine.Config
	v := viper.New()
	v.SetConfigType(format)

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	logrus.Debugf("flow script：%+v", cfg)
	return &cfg, nil
}

func ParserFile(filename string, format string) (*entity.Flow, error) {
	var flow *entity.Flow
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
	switch format {
	case "flowgram":
		flow, err = ParseFlowgram(data)
		if err != nil {
			logrus.Error("parse flowgram error:", err)
			return nil, fmt.Errorf("parse flowgram error, %s", err)
		}
		break
	default:
		logrus.Error("unknown format:", format)
		return nil, fmt.Errorf("unknown format: %s", format)
	}
	return flow, nil
}
