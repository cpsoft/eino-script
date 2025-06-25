package provider

import (
	"eino-script/engine"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (p *DataProvider) GetSessionList(flowId uint) (*[]Session, error) {
	var sessions []Session
	result := p.db.Where(&Session{FlowId: flowId}).Order("created_at DESC").Find(&sessions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询会话列表失败. %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return &sessions, nil
	}
	logrus.Debug("sessions", sessions)
	return &sessions, nil
}

func (p *DataProvider) NewSession(flowId uint, name string) (*engine.Session, error) {
	session := engine.NewSession(flowId)
	p.db.Save(&Session{SessionId: session.Id(), FlowId: flowId, Name: name})
	return session, nil
}

func (p *DataProvider) AddMessage(
	session *engine.Session,
	roleType schema.RoleType,
	content string) (*engine.Session, error) {
	session.AddMessage(roleType, content)
	p.db.Save(&SessionMessage{
		SessionId: session.Id(),
		Role:      roleType,
		Content:   content,
	})
	return session, nil
}

func (p *DataProvider) UpdateSession(id string, name string) error {
	result := p.db.Model(&Session{}).Where("session_id = ?", id).Update("name", name)
	if result.Error != nil {
		return fmt.Errorf("更新失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("更新记录不存在")
	}
	return nil
}

func (p *DataProvider) DeleteSession(id string) error {
	result := p.db.Model(&SessionMessage{}).Where("session_id = ?", id).Delete(&SessionMessage{})
	if result.Error != nil {
		return fmt.Errorf("删除模型失败：%s", result.Error.Error())
	}
	result = p.db.Model(&Session{}).Where("session_id = ?", id).Delete(&Session{})
	if result.Error != nil {
		return fmt.Errorf("删除会话失败：%s", result.Error.Error())
	}
	return nil
}

func (p *DataProvider) GetSessionMessages(sessionId string) (*[]SessionMessage, error) {
	var sessionMessages = make([]SessionMessage, 0)
	result := p.db.Where("session_id = ?", sessionId).Find(&sessionMessages)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logrus.Debugf("没有找到会话(%s)对应的消息", sessionId)
			return &[]SessionMessage{}, nil
		}
	}
	if result.RowsAffected == 0 {
		logrus.Debugf("没有找到会话(%s)对应的消息", sessionId)
		return &[]SessionMessage{}, nil
	}
	logrus.Debug("sessionMessages", sessionMessages)
	return &sessionMessages, nil
}
