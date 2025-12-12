@echo off
echo Quick build for current platform...

REM Определяем архитектуру
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=amd64
) else (
    set ARCH=386
)

echo Building ai-bot-windows-%ARCH%.exe...
go build -ldflags="-s -w" -o ai-bot-windows-%ARCH%.exe .

if %errorlevel% equ 0 (
    echo.
    echo ✅ Build successful: ai-bot-windows-%ARCH%.exe
    echo.
    echo You can now run:
    echo   ai-bot-windows-%ARCH%.exe --config
    echo   ai-bot-windows-%ARCH%.exe
) else (
    echo.
    echo ❌ Build failed!
    exit /b 1
)