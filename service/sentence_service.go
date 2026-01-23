package service

import (
	"daily-quote/database"
	"daily-quote/model"
)

// 随机获取一句话
func RandomSentence() (*model.Sentence, error) {
	var sentence model.Sentence
	err := database.DB.
		Order("RANDOM()").
		Limit(1).
		Find(&sentence).Error

	return &sentence, err
}

// 创建句子
func CreateSentence(sentence *model.Sentence) error {
	err := database.DB.Create(sentence).Error
	return err
}

// 删除单个句子
func DeleteSentence(id int) error {
	// 判断是否存在
	var s model.Sentence
	if err := database.DB.First(&s, id).Error; err != nil {
		return err
	}
	err := database.DB.Delete(&s, id).Error
	return err
}

// 分页查询句子
func ListSentences(page int, pageSize int) (*model.PageResult[model.Sentence], error) {
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
