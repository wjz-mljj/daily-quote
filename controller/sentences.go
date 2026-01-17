package controller

import (
	"daily-quote/database"
	"daily-quote/model"
	"daily-quote/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

// 分页查询句子
func GetListSentences(c *gin.Context) {
	data, err := service.ListSentences(1, 200)
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
		Data:    data,
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

// 导出
func ExportSentences(c *gin.Context) {
	data, err := service.ExportSentences()
	if err != nil {
		c.JSON(200, model.Response{
			Code:    500,
			Message: "Error",
			Data:    nil,
		})
		return
	}
	f := excelize.NewFile()
	sheet := "Sheet1"

	f.SetCellValue(sheet, "A1", "ID")
	f.SetCellValue(sheet, "B1", "内容")

	for i, s := range data {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), s.ID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), s.Content)
	}

	// 设置响应头（关键）
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=sentences.xlsx")

	_ = f.Write(c.Writer)
}
