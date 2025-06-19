package provider

import (
	"eino-script/engine/types"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// SaveFlow 创建或更新工作流记录，并返回结果
func (p *DataProvider) SaveFlow(info types.FlowInfo) (uint, error) {
	var resultID uint
	flow := Flow{
		Name:   info.Name,
		Script: info.Script,
	}
	flow.ID = info.ID
	err := p.db.Transaction(func(tx *gorm.DB) error {
		// 检查记录是否存在
		var existingFlow Flow
		if flow.ID != 0 {
			if err := tx.First(&existingFlow, flow.ID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// 记录不存在，视为创建操作
					result := tx.Create(&flow)
					if result.Error != nil {
						return fmt.Errorf("创建失败: %w", result.Error)
					}
					resultID = flow.ID
					return nil
				}
				return fmt.Errorf("查询失败: %w", err)
			}
		}

		// 更新现有记录
		if flow.ID == 0 || existingFlow.ID == 0 {
			// 创建新记录
			result := tx.Create(&flow)
			if result.Error != nil {
				return fmt.Errorf("创建失败: %w", result.Error)
			}
			resultID = flow.ID
		} else {
			// 更新现有记录
			result := tx.Model(&existingFlow).Updates(flow)
			if result.Error != nil {
				return fmt.Errorf("更新失败: %w", result.Error)
			}
			resultID = existingFlow.ID
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return resultID, nil
}

func (p *DataProvider) DeleteFlow(id uint) {
	p.db.Delete(&Flow{}, id)
}

func (p *DataProvider) GetFlowList() ([]Flow, error) {
	var flows []Flow
	result := p.db.Find(&flows)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Flow{}, nil
		}
		return nil, fmt.Errorf("查询工作流列表失败. %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return []Flow{}, nil
	}
	return flows, nil
}

func (p *DataProvider) GetFlow(id uint) (*Flow, error) {
	flow := Flow{}
	flow.ID = id
	result := p.db.First(&flow)
	if result.Error != nil {
		return nil, fmt.Errorf("查询工作流列表失败. %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("未找到工作流。%d", id)
	}
	return &flow, nil
}
