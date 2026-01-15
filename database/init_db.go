package database

import (
	"embed"
	"io"
	"os"
	"path/filepath"
)

// 从 embed 拷贝 SQLite 到本地
func InitDBFromEmbed(fs embed.FS) error {
	dbPath := "data/app.db"

	// 1️⃣ 如果数据库已存在，直接用
	if _, err := os.Stat(dbPath); err == nil {
		return nil
	}

	// 2️⃣ 创建目录
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return err
	}

	// 3️⃣ 从 embed 读取
	src, err := fs.Open("assets/app.db")
	if err != nil {
		return err
	}
	defer src.Close()

	// 4️⃣ 写入磁盘
	dst, err := os.Create(dbPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
