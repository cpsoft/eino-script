package engine

import "github.com/cloudwego/eino/schema"

type Session struct {
	messages []*schema.Message
}

func NewSession() *Session {
	messages := make([]*schema.Message, 0)
	return &Session{
		messages: messages,
	}
}

func (s *Session) Close() {
	return
}

func (s *Session) GetMessages() []*schema.Message {
	if len(s.messages) == 0 {
		return []*schema.Message{}
	}
	return s.messages
}

func (s *Session) AddMessage(roleType schema.RoleType, content string) {
	s.messages = append(s.messages, &schema.Message{
		Role:    roleType,
		Content: content,
	})
}
