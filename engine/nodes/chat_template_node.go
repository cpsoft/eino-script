package nodes

import (
	"eino-script/types"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

func CreateChatTemplateNode(cfg *types.NodeCfg) (prompt.ChatTemplate, error) {
	logrus.Infof("CreateChatTemplateNode: %+v", *cfg)

	var messagesTemplate = make([]schema.MessagesTemplate, 0)

	systemMsg, ok := cfg.Attrs["system_message"].(string)
	if !ok {
		return nil, fmt.Errorf("system_message not found in attrs")
	}
	messagesTemplate = append(messagesTemplate, schema.SystemMessage(systemMsg))

	history, ok := cfg.Attrs["history"].(bool)
	if ok && history {
		messagesTemplate = append(messagesTemplate, schema.MessagesPlaceholder("chat_history", true))
	}

	messagesTemplate = append(messagesTemplate, schema.UserMessage("问题：{message}?"))

	pt := prompt.FromMessages(
		schema.FString,
		messagesTemplate...,
	)

	return pt, nil
}
