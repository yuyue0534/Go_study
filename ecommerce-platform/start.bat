@echo off
chcp 65001 >nul

echo ==========================================
echo       优选商城 - 电商平台启动脚本
echo ==========================================
echo.

:: 检查Go是否安装
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 错误: Go语言环境未安装
    echo.
    echo 请先安装Go语言环境：
    echo   下载地址: https://golang.org/dl/
    echo.
    pause
    exit /b 1
)

for /f "tokens=*" %%i in ('go version') do set GO_VERSION=%%i
echo ✅ %GO_VERSION%
echo.

:: 进入项目目录
cd /d "%~dp0"

:: 下载依赖
echo 📦 正在下载依赖...
go mod tidy

if %errorlevel% neq 0 (
    echo ❌ 依赖下载失败
    pause
    exit /b 1
)

echo ✅ 依赖下载完成
echo.

:: 编译项目
echo 🔨 正在编译项目...
go build -o ecommerce-server.exe ./cmd/main.go

if %errorlevel% neq 0 (
    echo ❌ 编译失败
    pause
    exit /b 1
)

echo ✅ 编译完成
echo.

:: 启动服务
echo 🚀 正在启动服务...
echo.
ecommerce-server.exe
