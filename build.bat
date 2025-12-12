@echo off
echo Building AI Bot releases for multiple platforms...
echo.

REM Создаем папку для релизов
if not exist "releases" mkdir releases

REM Очищаем старые релизы
del /Q releases\* 2>nul

echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o releases\ai-bot-windows-amd64.exe .
if %errorlevel% neq 0 (
    echo Error building Windows amd64 version
    exit /b 1
)

echo Building for Windows (386)...
set GOOS=windows
set GOARCH=386
go build -ldflags="-s -w" -o releases\ai-bot-windows-386.exe .
if %errorlevel% neq 0 (
    echo Error building Windows 386 version
    exit /b 1
)

echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o releases\ai-bot-linux-amd64 .
if %errorlevel% neq 0 (
    echo Error building Linux amd64 version
    exit /b 1
)

echo Building for Linux (386)...
set GOOS=linux
set GOARCH=386
go build -ldflags="-s -w" -o releases\ai-bot-linux-386 .
if %errorlevel% neq 0 (
    echo Error building Linux 386 version
    exit /b 1
)

echo Building for Linux (arm64)...
set GOOS=linux
set GOARCH=arm64
go build -ldflags="-s -w" -o releases\ai-bot-linux-arm64 .
if %errorlevel% neq 0 (
    echo Error building Linux arm64 version
    exit /b 1
)

echo Building for FreeBSD (amd64)...
set GOOS=freebsd
set GOARCH=amd64
go build -ldflags="-s -w" -o releases\ai-bot-freebsd-amd64 .
if %errorlevel% neq 0 (
    echo Error building FreeBSD amd64 version
    exit /b 1
)

echo Building for FreeBSD (386)...
set GOOS=freebsd
set GOARCH=386
go build -ldflags="-s -w" -o releases\ai-bot-freebsd-386 .
if %errorlevel% neq 0 (
    echo Error building FreeBSD 386 version
    exit /b 1
)

echo Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o releases\ai-bot-darwin-amd64 .
if %errorlevel% neq 0 (
    echo Error building macOS amd64 version
    exit /b 1
)

echo Building for macOS (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o releases\ai-bot-darwin-arm64 .
if %errorlevel% neq 0 (
    echo Error building macOS arm64 version
    exit /b 1
)

echo.
echo All builds completed successfully!
echo.
echo Built files:
dir releases /b
echo.
echo Ready for GitHub release!