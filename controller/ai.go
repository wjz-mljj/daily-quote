package controller

import (
	"daily-quote/model"
	"daily-quote/service"
	"encoding/json"
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

type OllamaDeleteRequest struct { // 响应结构体
	ModelNmae string `json:"modelName"`
}

// 删除指定模型
func OllamaDeleteModele(c *gin.Context) {
	var params OllamaDeleteRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"details": err.Error(),
		})
		return
	}
	print(params.ModelNmae)
	reqs, err := service.Ollama_delete_model(params.ModelNmae)
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

// SSE 错误处理
func sseError(w gin.ResponseWriter, msg string) {
	data, _ := json.Marshal(gin.H{
		"status":  "error",
		"message": msg,
	})
	fmt.Fprintf(w, "data: %s\n\n", data)
	w.Flush()
}

// 拉取指定模型
func OllamaPullModel(c *gin.Context) {
	modelName := c.Query("modelName")
	if modelName == "" {
		c.JSON(400, gin.H{"error": "model name required"})
		return
	}
	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // nginx 下必须
	w.Flush()
	ch, err := service.Ollama_pull_model(modelName)
	if err != nil {
		sseError(w, err.Error())
		return
	}
	for progress := range ch {
		data, _ := json.Marshal(progress)

		fmt.Fprintf(w, "data: %s\n\n", data)
		w.Flush()

		if progress.Status == "success" {
			return
		}
	}

}
