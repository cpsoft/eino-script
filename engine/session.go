package engine

import (
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

type Session struct {
	id       string
	flowId   uint
	messages []*schema.Message
}

func NewSession(flowId uint) *Session {
	messages := make([]*schema.Message, 0)
	id, _ := uuid.NewUUID()
	return &Session{
		id:       id.String(),
		flowId:   flowId,
		messages: messages,
	}
}

func (s *Session) Close() {
	return
}

func (s *Session) AddMessage(roleType schema.RoleType, content string) {
	s.messages = append(s.messages, &schema.Message{
		Role:    roleType,
		Content: content,
	})
}

func (s *Session) Id() string {
	return s.id
}

func (s *Session) FlowId() uint {
	return s.flowId
}

func (s *Session) GetMessages() []*schema.Message {
	if len(s.messages) == 0 {
		return []*schema.Message{}
	}
	return s.messages
}
