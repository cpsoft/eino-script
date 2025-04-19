package engine

import (
	"eino-script/parser"
)

func (e *Engine) AddEdge(cfg *parser.EdgeCfg) error {
	_ = e.g.AddEdge(cfg.Src, cfg.Dst)
	return nil
}
