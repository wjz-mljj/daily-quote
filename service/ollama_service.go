package service

import (
	"bufio"
	"bytes"
	"daily-quote/database"
	"daily-quote/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var base_url = "http://127.0.0.1:11434"

// 使用 Ollama 生成文本分析结果
func OllamaGenerate(model_str string, sentence string, analysisType string, sentenceId uint) (interface{}, error, bool) {
	tpl := model.GetPromptTemplate(analysisType)
	// 渲染 Prompt （调用结构体的方法）
	prompt := tpl.Render(sentence)
	payload := model.OllamaGenerateRequest{
		Model:  model_str,
		System: model.SYSTEM_STR,
		Prompt: prompt,
		Stream: false,
		Think:  false,
		Options: model.Options{
			Temperature: tpl.Temperature,
			TopP:        0.9,
			NumPredict:  1024,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err, false
	}

	resp, err := http.Post(base_url+"/api/generate", "application/json",
		bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err, false
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err, false
	}

	// 写入
	var s model.Sentence
	if err := database.DB.First(&s, sentenceId).Error; err != nil {
		return result, nil, true
	}

	// 更新 这里只更新 analysis_results 和 type 字段 https://gorm.io/zh_CN/docs/update.html
	updateData := map[string]interface{}{
		"analysis_results": fmt.Sprintf("%v", result["response"]),
		"type":             tpl.Title,
	}

	if err := database.DB.Model(&s).Updates(updateData).Error; err != nil {
		return result, nil, true
	}

	// Markdown 转 HTML
	respStr, ok := result["response"].(string)
	if !ok {
		respStr = fmt.Sprintf("%v", result["response"])
	}
	var htmlStr string
	htmlStr, err = model.MarkdownToHTML(respStr)
	if err != nil {
		log.Println("转换失败:", err)
		htmlStr = respStr
	}
	result["response"] = htmlStr

	return result, nil, false
}

// 列出 Ollama 可用模型
func OllamaListModels() (interface{}, error) {
	resp, err := http.Get(base_url + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// 删除指定模型
func Ollama_delete_model(modelName string) (interface{}, error) {
	params := map[string]string{
		"model": modelName,
	}
	// 序列化为JSON
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", base_url+"/api/delete", bytes.NewBuffer(paramsBytes))
	if err != nil {
		return nil, err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return "success", nil

}

// 拉取指定模型
func Ollama_pull_model(modelName string) (<-chan model.PullProgress, error) {
	payload := map[string]any{
		"model":    modelName,
		"insecure": true,
		"stream":   true,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", base_url+"/api/pull", bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{Timeout: 0}).Do(req)
	if err != nil {
		return nil, err
	}
	// 创建一个通道用于传输拉取进度 使用协程读取响应体并发送进度到通道
	ch := make(chan model.PullProgress)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		for scanner.Scan() {
			var p model.PullProgress
			if err := json.Unmarshal(scanner.Bytes(), &p); err != nil {
				continue
			}

			ch <- p

			if p.Status == "success" {
				return
			}
		}
	}()

	return ch, nil
}
