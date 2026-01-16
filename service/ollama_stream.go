package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

var base_url = "http://127.0.0.1:11434"

type Options struct {
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	NumPredict  int     `json:"num_predict"`
}

type OllamaGenerateRequest struct {
	Model   string  `json:"model"`
	System  string  `json:"system"`
	Prompt  string  `json:"prompt"`
	Stream  bool    `json:"stream"`
	Think   bool    `json:"think"`
	Options Options `json:"options"`
}

var PROMPT_TEMPLATES = map[string]string{
	"semantic": `分析方向：语义分析。请从字面含义和隐含含义两个层面分析这句话。需要分析的句子：{sentence}`,
	"emotion":  `分析方向：情感分析。请判断这句话的情感倾向，并说明判断依据。需要分析的句子：{sentence}`,
	"humor":    `分析方向：幽默分析。请分析这句话是否具有幽默、讽刺或反讽效果，并说明原因。需要分析的句子：{sentence}`,
	"intent":   `分析方向：意图分析。请判断说话者的真实意图或目的。需要分析的句子：{sentence}`,
	"tone":     `分析方向：语气分析。请分析这句话的语气特征。需要分析的句子：{sentence}`,
}

func buildPrompt(analysisType, sentence string) (string, error) {
	template, exists := PROMPT_TEMPLATES[analysisType]
	if !exists {
		template = PROMPT_TEMPLATES["semantic"]
	}

	// 替换占位符
	result := strings.Replace(template, "{sentence}", sentence, -1)
	return result, nil
}

// 使用 Ollama 生成文本分析结果
func OllamaGenerate(model_str string, sentence string, analysisType string) (interface{}, error) {
	system_str := `
		你是一个专业的文本分析助手。
        你需要根据用户指定的【分析方向】，对给定句子进行分析。
        分析必须客观、清晰、有条理，禁止编造不存在的信息。
        输出请使用清晰的小标题结构。
	`
	prompt, _ := buildPrompt("semantic", sentence)

	payload := OllamaGenerateRequest{
		Model:  model_str,
		System: system_str,
		Prompt: prompt,
		Stream: false,
		Think:  false,
		Options: Options{
			Temperature: 0.2,
			TopP:        0.9,
			NumPredict:  1024,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(base_url+"/api/generate", "application/json",
		bytes.NewBuffer(payloadBytes))
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
	print(string(paramsBytes))
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

type PullProgress struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// 拉取指定模型
func Ollama_pull_model(modelName string) (<-chan PullProgress, error) {
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
	ch := make(chan PullProgress)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		for scanner.Scan() {
			var p PullProgress
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
