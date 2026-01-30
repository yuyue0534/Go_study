#!/bin/bash

# 电商平台启动脚本

echo "=========================================="
echo "      优选商城 - 电商平台启动脚本"
echo "=========================================="
echo ""

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: Go语言环境未安装"
    echo ""
    echo "请先安装Go语言环境："
    echo "  - Windows: https://golang.org/dl/"
    echo "  - macOS: brew install go"
    echo "  - Linux: sudo apt install golang-go"
    echo ""
    exit 1
fi

echo "✅ Go版本: $(go version)"
echo ""

# 进入项目目录
cd "$(dirname "$0")"

# 下载依赖
echo "📦 正在下载依赖..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ 依赖下载失败"
    exit 1
fi

echo "✅ 依赖下载完成"
echo ""

# 编译项目
echo "🔨 正在编译项目..."
go build -o ecommerce-server ./cmd/main.go

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译完成"
echo ""

# 启动服务
echo "🚀 正在启动服务..."
echo ""
./ecommerce-server
