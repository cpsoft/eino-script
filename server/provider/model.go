package provider

import (
	"eino-script/engine/types"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func (p *DataProvider) GetModelList() ([]Model, error) {
	var models []Model
	result := p.db.Find(&models)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Model{}, nil
		}
		return nil, fmt.Errorf("查询工作流列表失败. %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return []Model{}, nil
	}
	return models, nil
}

func (p *DataProvider) SaveModel(info *types.ModelInfo) error {
	model := Model{
		Name:             info.Name,
		ModelType:        info.ModelType,
		ModelName:        info.ModelName,
		ApiKey:           info.ApiKey,
		ApiUrl:           info.ApiUrl,
		MaxContextLength: info.MaxContextLength,
		StreamingEnabled: info.StreamingEnabled,
	}
	model.ID = info.ID

	result := p.db.Save(&model)
	if result.Error != nil {
		return fmt.Errorf("插入模型失败：%s", result.Error.Error())
	}
	return nil
}

func (p *DataProvider) DeleteModel(id uint) error {
	result := p.db.Delete(&Model{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除模型失败：%s", result.Error.Error())
	}
	return nil
}

func (p *DataProvider) GetModel(id uint) (*Model, error) {
	model := Model{}
	model.ID = id
	result := p.db.First(&model)
	if result.Error != nil {
		return nil, fmt.Errorf("查询工作流列表失败. %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("未找到工作流。%d", id)
	}
	return &model, nil
}
