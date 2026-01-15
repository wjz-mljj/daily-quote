package controller

import (
	"a-sentence/model"
	"a-sentence/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ModelsList(c *gin.Context) {
	reqs, err := service.OllamaListModels()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "error",
		})
		return
	}
	c.JSON(200, model.Response{
		Code:    200,
		Message: "success",
		Data:    reqs,
	})
}

type OllamaRequest struct {
	Model        string `json:"model"`
	Sentence     string `json:"sentence"`
	AnalysisType string `json:"analysis_type"`
}

func OllamaGenerateRequest(c *gin.Context) {
	var params OllamaRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"details": err.Error(),
		})
		return
	}
	fmt.Println(params.Model)
	reqs, err := service.OllamaGenerate(params.Model, params.Sentence, params.AnalysisType)

	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "error",
		})
		return
	}

	c.JSON(200, model.Response{
		Code:    200,
		Message: "success",
		Data:    reqs,
	})
}
