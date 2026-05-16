@echo off
setlocal
cd /d %~dp0\..
echo [1/2] Downloading dependencies...
go mod tidy
if errorlevel 1 (
  echo Dependency download failed.
  exit /b 1
)
echo [2/2] Starting student system...
go run .\cmd\server
endlocal
