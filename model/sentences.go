package model

import "gorm.io/gorm"

// gorm.Model 的定义：https://gorm.io/zh_CN/docs/models.html

// Sentence 句子模型
type Sentence struct {
	gorm.Model
	Content         string  `json:"content"`
	Type            *string `json:"type"`             // 句子分析类型， *表示可选字段
	AnalysisResults *string `json:"analysis_results"` // JSON string of analysis results
}

type PageResult[T any] struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
	List     []T   `json:"list"`
}
