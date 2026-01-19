package database

import (
	"daily-quote/model"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// https://gorm.io/zh_CN/docs/advanced_query.html
var DB *gorm.DB

func InitSQLite() {
	var err error

	DB, err = gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "data/app.db",
	}, &gorm.Config{})

	if err != nil {
		log.Fatal("SQLite 连接失败:", err)
	}
	DB.AutoMigrate(&model.Sentence{})

	log.Println("SQLite 连接成功")
}
