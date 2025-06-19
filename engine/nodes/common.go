package nodes

import (
	"eino-script/engine/types"
	"errors"
	"github.com/cloudwego/eino/components/prompt"
)

type NodesFactroyData struct {
}

var chattemplate_functions = map[types.NodeType]func(cfg *types.NodeCfg) (prompt.ChatTemplate, error){
	"ChatTemplate": CreateChatTemplateNode,
}

func CreateChatTemplateByFactroy(cfg *types.NodeCfg) (prompt.ChatTemplate, error) {
	f, ok := chattemplate_functions[cfg.Type]
	if !ok {
		return nil, errors.New("Node not found")
	}
	return f(cfg)
}
