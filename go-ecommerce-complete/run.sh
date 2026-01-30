#!/bin/bash

echo "======================================"
echo "     Go电商平台 - 启动脚本"
echo "======================================"
echo ""

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误：未检测到Go环境"
    echo "请访问 https://golang.org/dl/ 安装Go"
    exit 1
fi

echo "✓ Go版本: $(go version)"
echo ""

# 安装依赖
echo "正在安装依赖..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ 依赖安装失败"
    exit 1
fi

echo "✓ 依赖安装完成"
echo ""

# 启动服务器
echo "正在启动服务器..."
echo ""
go run cmd/server/main.go
