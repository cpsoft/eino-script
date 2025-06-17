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
	data, ok := cfg.Attrs["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in attrs")
	}

	logrus.Debug("attrs:", data)

	systemMsg, ok := data["systemmessage"].(string)
	if !ok {
		return nil, fmt.Errorf("systemMessage not found in attrs")
	}
	messagesTemplate = append(messagesTemplate, schema.SystemMessage(systemMsg))

	history, ok := data["history"].(bool)
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
