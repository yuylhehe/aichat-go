# AI Chat 部署文档

本文档介绍了如何将 AI Chat 应用打包并部署到 Linux 服务器。

## 1. 打包

在 macOS 开发环境中，执行以下命令进行交叉编译：

```bash
chmod +x build_linux.sh
./build_linux.sh
```

执行成功后，会在当前目录下生成 `aichat-linux-amd64` 可执行文件。该文件已经包含了前端静态资源，是一个独立的二进制文件。

## 2. 部署到 Linux 服务器

### 2.1 准备工作

将以下文件上传到服务器的同一目录下：

1. `aichat-linux-amd64` (可执行文件)
2. `.env` (配置文件)

### 2.2 运行

在服务器上，首先赋予可执行权限：

```bash
chmod +x aichat-linux-amd64
```

然后启动应用：

```bash
./aichat-linux-amd64
```

### 2.3 后台运行 (使用 Systemd)

建议使用 Systemd 来管理服务。

1. 创建服务文件 `/etc/systemd/system/aichat.service`：

```ini
[Unit]
Description=AI Chat Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/path/to/your/app
ExecStart=/path/to/your/app/aichat-linux-amd64
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

请将 `/path/to/your/app` 替换为实际的部署路径。

2. 启动服务：

```bash
systemctl daemon-reload
systemctl enable aichat
systemctl start aichat
```

3. 查看状态：

```bash
systemctl status aichat
```

## 3. 注意事项

- **配置文件**: 确保 `.env` 文件中的配置项（如数据库连接、端口等）在服务器环境下是正确的。
- **静态资源**: 前端资源已经嵌入到二进制文件中，无需单独上传 `public` 目录。如果需要更新前端，需要重新打包后端。
- **数据库**: 确保服务器能够访问到配置的数据库。
