#!/bin/bash

# 设置目标操作系统和架构
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# 输出文件名
OUTPUT_NAME="aichat-linux-amd64"

echo "Building for Linux/AMD64..."
go build -o $OUTPUT_NAME main.go

if [ $? -eq 0 ]; then
    echo "Build successful! Output file: $OUTPUT_NAME"
else
    echo "Build failed!"
    exit 1
fi
