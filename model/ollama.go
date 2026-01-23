package model

import (
	"strings"
)

type OllamaRequest struct {
	Model        string `json:"model"`
	Sentence     string `json:"sentence"`
	AnalysisType string `json:"analysis_type"`
	SentenceId   uint   `json:"sentence_id"`
}

type OllamaDeleteRequest struct { // 响应结构体
	ModelNmae string `json:"modelName"`
}

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

// 渲染模板(PromptTemplate结构体 增加 Render 方法)
func (pt PromptTemplate) Render(sentence string) string {
	return strings.ReplaceAll(pt.Prompt, "{sentence}", sentence)
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
