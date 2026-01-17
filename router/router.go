package router

import (
	"daily-quote/controller"
	"daily-quote/middleware"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(webFS embed.FS) *gin.Engine {
	// gin.ReleaseMode gin.DebugMode
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(middleware.Recovery())

	isDev := gin.Mode() == gin.DebugMode
	log.Println("运行模式:", isDev)
	if isDev {
		r.Static("/static", "./web")

		r.GET("/", func(c *gin.Context) {
			c.File("./web/index.html")
		})
	} else {
		// 静态文件服务
		sub, _ := fs.Sub(webFS, "web")
		r.StaticFS("/static", http.FS(sub))
		r.GET("/", func(c *gin.Context) {
			data, err := webFS.ReadFile("web/index.html")
			if err != nil {
				c.String(500, "index.html not found")
				return
			}
			c.Data(200, "text/html; charset=utf-8", data)
		})
	}

	api := r.Group("/api")
	{
		api.GET("/ping", controller.Ping)
		api.POST("/add_sentence", controller.CreateSentence)
		api.GET("/single_sentence", controller.GetRandomSentence)
		api.GET("/sentence_list", controller.GetListSentences)
		api.GET("/ollama_models", controller.ModelsList)
		api.POST("/ollama_generate", controller.OllamaGenerateRequest)
		api.POST("/ollama_delete_model", controller.OllamaDeleteModele)
		api.GET("/pull/stream", controller.OllamaPullModel)
		r.GET("/export/sentences", controller.ExportSentences)
	}

	return r
}
