package model

import "gorm.io/gorm"

type Sentence struct {
	gorm.Model
	Content string `json:"content"`
}
