package controller

import (
	"daily-quote/model"
	"daily-quote/service"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// CreateSentence 创建句子
func CreateSentence(c *gin.Context) {
	var sentence model.Sentence
	// 必须参数绑定 参数是JSON格式
	if err := c.ShouldBindJSON(&sentence); err != nil {
		model.Fail(c, 400, "请求参数错误")
		return
	}
	err := service.CreateSentence(&sentence)
	if err != nil {
		model.Fail(c, 500, "创建句子失败")
		return
	}
	model.Success(c, "创建句子成功", sentence)
}

// DeleteSentence 删除句子
func DeleteSentence(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		model.Fail(c, 400, "请求参数错误")
		return
	}

	err = service.DeleteSentence(id)
	if err != nil {
		model.Fail(c, 500, "删除句子失败")
		return
	}
	model.Success(c, "删除句子成功", nil)
}

// 分页查询句子
type ListSentencesRequest struct { // 定义请求参数结构体 form标签用于绑定URL查询参数
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func GetListSentences(c *gin.Context) {
	var params ListSentencesRequest
	// 绑定 URL 查询参数
	if err := c.ShouldBindQuery(&params); err != nil {
		model.Fail(c, 400, "请求参数错误")
		return
	}
	fmt.Printf("page: %T, page: %d, pageSize: %T, pageSize: %d\n", params.Page, params.Page, params.PageSize, params.PageSize)
	data, err := service.ListSentences(params.Page, params.PageSize)
	if err != nil {
		model.Fail(c, 500, "Error")
		return
	}
	model.Success(c, "success", data)
}

// GetRandomSentence 获取单个句子
func GetRandomSentence(c *gin.Context) {
	sentence, err := service.RandomSentence()
	if err != nil {
		model.Fail(c, 500, "Error")
		return
	}
	model.Success(c, "success", sentence)
}

// 导出 excel
func ExportSentences(c *gin.Context) {
	data, err := service.ExportSentences()
	if err != nil {
		model.Fail(c, 500, "Error")
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
