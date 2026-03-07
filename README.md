# 每日一句 - 本地 Ollama 驱动的文本分析工具

一个基于 Golang + 本地 Ollama 构建的「每日一句」分析项目，支持本地存储、多维度文本分析，所有数据全程本地化，不上传、不出本机。

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8.svg)](https://golang.org/)
[![Ollama](https://img.shields.io/badge/Ollama-Local-7D64FF.svg)](https://ollama.com/)

## 🌟 功能特性

### 📚 本地名言管理
- 使用 SQLite 数据库本地化存储所有数据
- 支持用户添加、删除、查询名言名句
- 数据仅保存在本地设备，不上传任何云端服务器

### 🤖 本地大模型分析
- 通过 Ollama 调用本地部署的大语言模型
- 支持下载模型、删除本地模型
- 无需依赖云端 API 接口
- 零隐私泄露风险、零使用成本、无网络依赖

### 🧠 多分析方向可选
同一句话可从不同维度进行深度理解和拆解：
- **文学表达分析**
- **逻辑与含义拆解**
- **情绪与语气判断**
- **沟通场景解读**
- **学习与思考角度延展**
### ⚠️ 说明
- 本项目不依赖云端 API
- 所有名言仅保存在本地
- 模型响应速度取决于本机性能和所使用的模型

## 🛠 技术栈
- **开发语言**：Golang
- **数据库**：SQLite（本地文件数据库）
- **大模型运行**：Ollama（本地部署）
- **分析策略**：基于不同分析方向的 Prompt Engineering

## 🚀 快速开始

### 环境准备
请确保本机已安装并准备好以下环境：
- Go ≥ 1.20
- Ollama（本地运行）

确认 Ollama 已正常运行：
```bash
ollama --version
```
### 获取项目代码
```bash
git clone https://gitee.com/w-jz/daily-quote.git
cd daily-quote
```
### 启动项目
```bash
go mod tidy
go run main.go
```
项目启动后：
- 会在本地创建 SQLite 数据库文件
- 不会上传任何数据
- 所有分析请求均通过 本地 Ollama 完成(需确保本地Ollama运行)
- 在浏览器中输入：http://127.0.0.1:8901

### 打包项目
> 打包命令：go build -o app.exe

> mac 打包windows命令：GOOS=windows GOARCH=amd64 go build -o app.exe

> mac 打包命令：GIN_MODE=release go build -o app main.go


## 📦 项目地址
Github: [https://github.com/wjz-mljj/daily-quote](https://github.com/wjz-mljj/daily-quote)
