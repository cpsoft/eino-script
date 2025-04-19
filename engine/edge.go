package engine

import (
	"eino-script/types"
	"github.com/sirupsen/logrus"
)

func (e *Engine) AddEdge(cfg *types.EdgeCfg) error {
	logrus.Infof("CreateEdge: %s -> %s", cfg.Src, cfg.Dst)
	err := e.g.AddEdge(cfg.Src, cfg.Dst)
	if err != nil {
		return err
	}
	return nil
}
