package provider

import (
	"eino-script/engine/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (p *DataProvider) GetMcpList() ([]Mcp, error) {
	var mcp []Mcp
	result := p.db.Find(&mcp)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Mcp{}, nil
		}
		return nil, fmt.Errorf("查询工作流列表失败. %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return []Mcp{}, nil
	}
	return mcp, nil
}

func (p *DataProvider) SaveMcp(info *types.McpInfo) (uint, error) {
	logrus.Debug("SaveMcp:", *info)
	var resultID uint

	raw, err := json.Marshal(info.McpCaps.Prompts)
	if err != nil {
		return 0, fmt.Errorf("MCP Prompts 数据错误.")
	}
	prompts := string(raw)
	logrus.Debug("prompts:", prompts)

	raw, err = json.Marshal(info.McpCaps.Tools)
	if err != nil {
		return 0, fmt.Errorf("MCP Tools 数据错误.")
	}
	tools := string(raw)

	raw, err = json.Marshal(info.McpCaps.Resources)
	if err != nil {
		return 0, fmt.Errorf("MCP Tools 数据错误.")
	}
	resources := string(raw)

	mcp := Mcp{
		Name:      info.Name,
		McpType:   info.McpType,
		Url:       info.Url,
		Prompts:   prompts,
		Tools:     tools,
		Resources: resources,
	}
	mcp.ID = info.ID

	logrus.Debug("MCP:", mcp)

	err = p.db.Transaction(func(tx *gorm.DB) error {
		// 检查记录是否存在
		var existingMcp Mcp
		if mcp.ID != 0 {
			if err := tx.First(&existingMcp, mcp.ID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// 记录不存在，视为创建操作
					result := tx.Create(&mcp)
					if result.Error != nil {
						return fmt.Errorf("创建失败: %w", result.Error)
					}
					resultID = mcp.ID
					return nil
				}
				return fmt.Errorf("查询失败: %w", err)
			}
		}

		// 更新现有记录
		if mcp.ID == 0 || existingMcp.ID == 0 {
			// 创建新记录
			result := tx.Create(&mcp)
			if result.Error != nil {
				return fmt.Errorf("创建失败: %w", result.Error)
			}
			resultID = mcp.ID
		} else {
			// 更新现有记录
			result := tx.Model(&existingMcp).Updates(mcp)
			if result.Error != nil {
				return fmt.Errorf("更新失败: %w", result.Error)
			}
			resultID = existingMcp.ID
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return resultID, nil
}

func (p *DataProvider) DeleteMcp(id uint) error {
	result := p.db.Delete(&Mcp{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除模型失败：%s", result.Error.Error())
	}
	return nil
}

func (p *DataProvider) GetMcp(id uint) (*Mcp, error) {
	mcp := Mcp{}
	mcp.ID = id
	result := p.db.First(&mcp)
	if result.Error != nil {
		return nil, fmt.Errorf("查询MCP失败. %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("未找到MCP。%d", id)
	}
	return &mcp, nil
}
