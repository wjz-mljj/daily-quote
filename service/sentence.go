package service

import (
	"daily-quote/database"
	"daily-quote/model"
)

func RandomSentence() (*model.Sentence, error) {
	var sentence model.Sentence
	err := database.DB.
		Order("RANDOM()").
		Limit(1).
		Find(&sentence).Error

	return &sentence, err
}

// 分页查询句子
func ListSentences(page, pageSize int) (*model.PageResult[model.Sentence], error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var total int64
	err := database.DB.
		Model(&model.Sentence{}).
		Count(&total).Error
	if err != nil {
		return nil, err
	}

	var list_a []model.Sentence
	err = database.DB.
		Model(&model.Sentence{}).
		Order("id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&list_a).Error
	if err != nil {
		return nil, err
	}

	return &model.PageResult[model.Sentence]{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		List:     list_a,
	}, nil
}

// 导出
func ExportSentences() ([]model.Sentence, error) {
	var list_a []model.Sentence
	err := database.DB.Find(&list_a).Error

	return list_a, err
}
