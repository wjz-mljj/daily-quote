package service

import (
	"a-sentence/database"
	"a-sentence/model"
)

func RandomSentence() (*model.Sentence, error) {
	var sentence model.Sentence
	err := database.DB.
		Order("RANDOM()").
		Limit(1).
		Find(&sentence).Error

	return &sentence, err
}
