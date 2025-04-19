package engine

import (
	"eino-script/parser"
	"github.com/sirupsen/logrus"
)

func (e *Engine) BindTools(modelName string, toolsName string) (string, error) {
	var err error
	m, ok := e.models[modelName]
	if !ok {
		return modelName, nil
	}
	toolsNode, ok := e.tools[toolsName]
	if !ok {
		return modelName, nil
	}
	logrus.Infof("BindTools: %s <- %s\n", modelName, toolsName)
	newModel, err := m.WithTools(toolsNode)
	if err != nil {
		return modelName, err
	}
	newName := modelName + "_" + "withTools"
	e.models[newName] = newModel
	err = e.g.AddChatModelNode(newName, newModel)
	if err != nil {
		return "", err
	}
	return newName, nil
}

func (e *Engine) AddEdge(cfg *parser.EdgeCfg) error {
	logrus.Infof("CreateEdge: %s -> %s\n", cfg.Src, cfg.Dst)
	name, err := e.BindTools(cfg.Src, cfg.Dst)
	if err != nil {
		return err
	}
	err = e.g.AddEdge(name, cfg.Dst)
	if err != nil {
		return err
	}
	return nil
}
