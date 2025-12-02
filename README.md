# AI Chat - High Performance Go AI Chat System

> 一个基于 Go (Gin) 和原生 JavaScript 构建的高性能流式 AI 对话系统。支持 OpenAI 格式接口，适配 GLM-4.6 等推理模型，提供极致的打字机体验。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)
![Gin](https://img.shields.io/badge/Gin-v1.10-00ADD8.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791.svg)

## ✨ 核心特性 (Features)

- **🚀 全链路流式响应 (End-to-End Streaming)**
  - 后端采用 Go 协程 + Channel 实现 Producer-Consumer 模式，前端使用 EventSource (SSE)，实现毫秒级首字延迟。
  - 支持连接状态感知，用户断开连接时自动停止生成并保存已生成内容。

- **🧠 深度推理支持 (Reasoning Support)**
  - 完美适配 GLM-4.6 等具备推理能力的模型。
  - 专门设计的 UI 支持折叠/展开“思考过程” (Chain of Thought)，让用户既能看到结果也能理解逻辑。

- **🛡️ 生产级架构**
  - **分层设计**: Handler-Service-Repository 清晰分层，易于维护和扩展。
  - **安全可靠**: 内置 JWT 认证、CORS 跨域配置、HTTP 代理支持。
  - **运维友好**: 支持 Systemd 托管，提供 Linux 交叉编译脚本，单文件部署 (Single Binary)。

- **💻 轻量级前端**
  - 无需 Webpack/Vite 构建，采用原生 ES Modules 开发。
  - Tailwind CSS 加持，界面简洁美观，加载速度极快。

## 🛠️ 技术栈 (Tech Stack)

- **Backend**: Go 1.22+, Gin, GORM, JWT
- **Database**: PostgreSQL
- **Frontend**: Native JavaScript (ESM), Tailwind CSS, EventSource
- **Infrastructure**: Docker (Optional), Systemd, Nginx

## 🚀 快速开始 (Quick Start)

### 1. 环境准备
- Go 1.22+
- PostgreSQL 14+

### 2. 克隆项目
```bash
git clone https://github.com/your-username/aichat-go.git
cd aichat-go
```

### 3. 配置环境变量
复制 `.env.example` 为 `.env` 并填入你的配置：

```bash
# Database
DATABASE_URL="postgresql://postgres:password@host:port/database"

# AI Service (OpenAI Compatible)
AI_API_KEY=sk-xxxxxx
AI_BASE_URL=https://api.deepseek.com/v1
AI_MODEL=deepseek-chat
```

### 4. 运行项目

**开发模式:**
```bash
go run main.go
```

**编译运行:**
```bash
go build -o aichat
./aichat
```

访问浏览器: `http://localhost:8080`

## 📦 部署指南 (Deployment)

### 交叉编译 (macOS -> Linux)
项目提供了便捷的构建脚本：
```bash
chmod +x build_linux.sh
./build_linux.sh
```
这将生成 `aichat-linux-amd64` 二进制文件，直接上传到服务器即可运行。

### Systemd 托管
建议在 Linux 服务器上使用 Systemd 管理进程：

```ini
# /etc/systemd/system/aichat.service
[Unit]
Description=AI Chat Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/www/aichat
ExecStart=/www/aichat/aichat-linux-amd64
Restart=always

[Install]
WantedBy=multi-user.target
```

## 📂 目录结构 (Structure)

```
├── assets/             # 静态资源 (嵌入二进制)
│   ├── public/         # 前端代码 (HTML/JS/CSS)
│   └── assets.go       # embed 声明
├── config/             # 配置加载
├── internal/           # 业务逻辑
│   ├── handler/        # HTTP 接口层
│   ├── service/        # 业务逻辑层
│   ├── repository/     # 数据访问层
│   ├── model/          # 数据库模型
│   └── middleware/     # Gin 中间件
├── build_linux.sh      # 构建脚本
└── main.go             # 入口文件
```

## 🤝 贡献 (Contributing)
欢迎提交 Issue 和 Pull Request！

## 📄 许可证 (License)
MIT License
