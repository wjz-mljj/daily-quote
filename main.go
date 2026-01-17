package main

import (
	"daily-quote/database"
	"daily-quote/router"
	"embed"
	"log"
)

//go:embed web/*
var webFS embed.FS

//go:embed assets/app.db
var dbFS embed.FS

func main() {
	if err := database.InitDBFromEmbed(dbFS); err != nil {
		log.Fatal(err)
	}
	database.InitSQLite()
	r := router.InitRouter(webFS)

	addr := "127.0.0.1:8901"
	log.Println("服务启动成功：http://" + addr)
	if err := r.Run(":8901"); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}
