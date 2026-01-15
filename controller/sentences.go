package controller

import (
	"a-sentence/database"
	"a-sentence/model"
	"a-sentence/service"

	"github.com/gin-gonic/gin"
)

// CreateSentence 创建句子
func CreateSentence(c *gin.Context) {
	var sentence model.Sentence

	if err := c.ShouldBindJSON(&sentence); err != nil {
		c.JSON(200, model.Response{
			Code:    400,
			Message: "请求参数错误",
			Data:    nil,
		})
		return
	}

	if err := database.DB.Create(&sentence).Error; err != nil {
		c.JSON(200, model.Response{
			Code:    500,
			Message: "创建句子失败",
			Data:    nil,
		})
		return
	}
	c.JSON(200, model.Response{
		Code:    200,
		Message: "创建句子成功",
		Data:    sentence,
	})
}

// ListUsers 列出所有句子
func ListSentence(c *gin.Context) {
	var sentences []model.Sentence
	database.DB.Find(&sentences)

	c.JSON(200, gin.H{
		"code": 0,
		"data": sentences,
	})
}

// GetRandomSentence 获取单个句子
func GetRandomSentence(c *gin.Context) {
	sentence, err := service.RandomSentence()
	if err != nil {
		c.JSON(200, model.Response{
			Code:    500,
			Message: "Error",
			Data:    nil,
		})
		return
	}
	c.JSON(200, model.Response{
		Code:    200,
		Message: "success",
		Data:    sentence,
	})
}
