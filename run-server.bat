@echo off
title GoTemplate API Server
cd /d "%~dp0"

echo ========================================
echo   GoTemplate API Server
echo ========================================
echo.

echo Building application...
go build -o tmp/main.exe .

if errorlevel 1 (
    echo ❌ Build failed!
    pause
    exit /b 1
)

echo ✅ Build successful!
echo.
echo Starting server...
echo.

tmp\main.exe

echo.
echo ========================================
echo Server stopped.
pause