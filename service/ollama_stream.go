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

type PullProgress struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

type PromptTemplate struct {
	Prompt      string  `json:"prompt"`
	Title       string  `json:"title"`
	Temperature float64 `json:"temperature"`
}

var SYSTEM_STR = `
	你是一个“语言与沟通分析型 AI”，专门对用户输入的一句话进行多角度、结构化分析。
	你的目标不是简单解释字面意思，而是：
	- 帮助用户理解这句话“怎么说的”
	- “为什么这么说”
	- “在真实沟通中会被如何理解”
	- 以及“如果换一种说法会发生什么变化”

	你需要遵守以下原则：
	1. 所有分析必须基于原句本身，不进行无依据的心理揣测
	2. 语言清晰、逻辑分层，避免空泛、套话
	3. 以“学习与提升表达能力”为导向，而非批评或说教
	4. 不使用过度学术化术语，普通人可以看懂
	5. 输出内容结构化，必要时使用分点说明

	你需要根据用户指定的【方向】，对给定句子进行分析。
`

var PROMPT_TEMPLATES_TWO = map[string]PromptTemplate{
	"literary": {
		Title:       "文学表达分析",
		Temperature: 0.4,
		Prompt: `
			方向：文学表达分析。
			请从“文学与语言表达”的角度分析这句话，包括但不限于：
			- 用词特点（是否口语化、抽象、形象）
			- 句式结构（简短 / 递进 / 留白等）
			- 是否存在修辞、隐含表达或风格倾向
			重点放在“表达方式本身”，而不是情绪或立场。
			需要分析的句子：{sentence}
		`,
	},
	"logic": {
		Title:       "逻辑与含义拆解",
		Temperature: 0.2,
		Prompt: `
			方向：逻辑与含义拆解。
			请对这句话进行逻辑与含义拆解，包括：
			- 表层意思（字面在说什么）
			- 隐含前提（默认成立的背景或假设）
			- 核心指向（真正想表达的重点）
			- 是否存在模糊、跳跃或多义理解空间
			需要分析的句子：{sentence}
		`,
	},
	"emotion": {
		Title:       "情绪与语气判断",
		Temperature: 0.4,
		Prompt: `
			方向：情绪与语气判断。
			请分析这句话中可能体现的情绪与语气特征，包括：
			- 情绪倾向（如冷静、不满、调侃、无奈等）
			- 语气强度（克制 / 明显 / 含蓄）
			- 情绪是直接表达还是通过措辞间接体现

			注意：只基于语言本身判断，不做人格或动机定性。
			需要分析的句子：{sentence}
		`,
	},
	"context": {
		Title:       "口语 / 沟通场景解读",
		Temperature: 0.5,
		Prompt: `
			方向：口语 / 沟通场景解读。
			请从真实沟通场景的角度解读这句话：
			- 这句话更像出现在什么场合（工作 / 亲密关系 / 社交等）
			- 听者可能会如何理解或产生何种感受
			- 这句话在沟通中起到的作用（拉近关系 / 划清界限 / 试探 / 回避等）
			需要分析的句子：{sentence}
		`,
	},
	"learning": {
		Title:       "学习与思考角度延展",
		Temperature: 0.7,
		Prompt: `
			方向：学习与思考角度延展。
			请从学习与思考的角度进行延展：
			- 这句话在表达或沟通上有什么值得学习的地方
			- 或者，它可能暴露了哪些常见的表达问题
			- 给用户 1–2 个“思考角度”，帮助其在今后类似表达中做出更有意识的选择
			请控制思考延展的范围，紧密围绕原句的表达方式与沟通影响，避免泛泛而谈或抽象价值判断。
			需要分析的句子：{sentence}
		`,
	},
}

// 获取模板
func GetPromptTemplate(templateKey string) PromptTemplate {
	tpl, exists := PROMPT_TEMPLATES_TWO[templateKey]
	if !exists {
		tplT, _ := PROMPT_TEMPLATES_TWO["literary"]
		tpl = tplT
	}
	return tpl
}

// 渲染模板
func (pt PromptTemplate) Render(sentence string) string {
	return strings.ReplaceAll(pt.Prompt, "{sentence}", sentence)
}

// 使用 Ollama 生成文本分析结果
func OllamaGenerate(model_str string, sentence string, analysisType string, sentenceId uint) (interface{}, error, bool) {
	tpl := GetPromptTemplate(analysisType)
	prompt := tpl.Render(sentence)
	payload := OllamaGenerateRequest{
		Model:  model_str,
		System: SYSTEM_STR,
		Prompt: prompt,
		Stream: false,
		Think:  false,
		Options: Options{
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
