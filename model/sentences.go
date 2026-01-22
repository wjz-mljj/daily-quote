package model

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"gorm.io/gorm"
)

// gorm.Model 的定义：https://gorm.io/zh_CN/docs/models.html

// Sentence 句子模型
type Sentence struct {
	gorm.Model
	Content         string `json:"content"`
	Type            string `json:"type" gorm:"default:''"` // 句子分析类型 默认空字符串
	AnalysisResults string `json:"analysis_results" gorm:"default:''"`
}

type PageResult[T any] struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
	List     []T   `json:"list"`
}

// MarkdownToPlainText 将Markdown格式的文本转换为纯文本
func MarkdownToPlainText(md string) string {
	source := []byte(md)

	parser := goldmark.New()
	reader := text.NewReader(source)
	doc := parser.Parser().Parse(reader)

	var buf bytes.Buffer

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch node := n.(type) {

		case *ast.Text:
			if entering {
				buf.Write(node.Segment.Value(source))
			}

		case *ast.Paragraph, *ast.Heading, *ast.ListItem:
			// 块级元素结束时换行
			if !entering {
				buf.WriteByte('\n')
			}
		}
		return ast.WalkContinue, nil
	})

	return buf.String()
}

// MarkdownToHTML 将Markdown格式的文本转换为HTML
func MarkdownToHTML(md string) (string, error) {
	var buf bytes.Buffer

	parser := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // 表格、列表、删除线
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // 单换行 -> <br>
			html.WithUnsafe(),    // 如果 LLM 返回 HTML（可选）
		),
	)

	err := parser.Convert([]byte(md), &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
