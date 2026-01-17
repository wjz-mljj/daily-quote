package model

import "gorm.io/gorm"

type Sentence struct {
	gorm.Model
	Content string `json:"content"`
}

type PageResult[T any] struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
	List     []T   `json:"list"`
}
