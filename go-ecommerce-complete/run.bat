@echo off
echo ======================================
echo      Go电商平台 - 启动脚本
echo ======================================
echo.

REM 检查Go环境
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [错误] 未检测到Go环境
    echo 请访问 https://golang.org/dl/ 安装Go
    pause
    exit /b 1
)

for /f "tokens=*" %%i in ('go version') do set GOVER=%%i
echo [成功] Go版本: %GOVER%
echo.

REM 安装依赖
echo 正在安装依赖...
go mod tidy

if %ERRORLEVEL% NEQ 0 (
    echo [错误] 依赖安装失败
    pause
    exit /b 1
)

echo [成功] 依赖安装完成
echo.

REM 启动服务器
echo 正在启动服务器...
echo.
go run cmd/server/main.go

pause
