@echo off
echo Building AI Bot complete releases with documentation...
echo.

REM Создаем папки
if not exist "releases" mkdir releases
if not exist "temp" mkdir temp

REM Очищаем старые релизы
del /Q releases\* 2>nul
rmdir /S /Q temp 2>nul
mkdir temp

echo Building Windows releases...

REM Windows 64-bit
echo Building Windows 64-bit...
mkdir temp\ai-bot-windows-amd64
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o temp\ai-bot-windows-amd64\ai-bot.exe .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-windows-amd64\ >nul 2>&1
copy USAGE.md temp\ai-bot-windows-amd64\ >nul 2>&1
copy .env.example temp\ai-bot-windows-amd64\ >nul 2>&1

REM Создаем install.bat
echo @echo off > temp\ai-bot-windows-amd64\install.bat
echo echo 🤖 AI Bot Installation >> temp\ai-bot-windows-amd64\install.bat
echo echo ====================== >> temp\ai-bot-windows-amd64\install.bat
echo echo. >> temp\ai-bot-windows-amd64\install.bat
echo echo ✅ AI Bot готов к использованию! >> temp\ai-bot-windows-amd64\install.bat
echo echo. >> temp\ai-bot-windows-amd64\install.bat
echo echo Быстрый старт: >> temp\ai-bot-windows-amd64\install.bat
echo echo   ai-bot.exe --config  # Настройка >> temp\ai-bot-windows-amd64\install.bat
echo echo   ai-bot.exe           # Запуск >> temp\ai-bot-windows-amd64\install.bat
echo echo. >> temp\ai-bot-windows-amd64\install.bat
echo pause >> temp\ai-bot-windows-amd64\install.bat

powershell -Command "Compress-Archive -Path 'temp\ai-bot-windows-amd64\*' -DestinationPath 'releases\ai-bot-windows-amd64.zip' -Force"
echo ✅ ai-bot-windows-amd64.zip

REM Windows 32-bit
echo Building Windows 32-bit...
mkdir temp\ai-bot-windows-386
set GOOS=windows
set GOARCH=386
go build -ldflags="-s -w" -o temp\ai-bot-windows-386\ai-bot.exe .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-windows-386\ >nul 2>&1
copy USAGE.md temp\ai-bot-windows-386\ >nul 2>&1
copy .env.example temp\ai-bot-windows-386\ >nul 2>&1
copy temp\ai-bot-windows-amd64\install.bat temp\ai-bot-windows-386\ >nul

powershell -Command "Compress-Archive -Path 'temp\ai-bot-windows-386\*' -DestinationPath 'releases\ai-bot-windows-386.zip' -Force"
echo ✅ ai-bot-windows-386.zip

echo.
echo Building Linux releases...

REM Linux 64-bit
echo Building Linux 64-bit...
mkdir temp\ai-bot-linux-amd64
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o temp\ai-bot-linux-amd64\ai-bot .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-linux-amd64\ >nul 2>&1
copy USAGE.md temp\ai-bot-linux-amd64\ >nul 2>&1
copy .env.example temp\ai-bot-linux-amd64\ >nul 2>&1

REM Создаем install.sh
echo #!/bin/bash > temp\ai-bot-linux-amd64\install.sh
echo # AI Bot Installation Script >> temp\ai-bot-linux-amd64\install.sh
echo. >> temp\ai-bot-linux-amd64\install.sh
echo echo "🤖 AI Bot Installation" >> temp\ai-bot-linux-amd64\install.sh
echo echo "======================" >> temp\ai-bot-linux-amd64\install.sh
echo echo "" >> temp\ai-bot-linux-amd64\install.sh
echo. >> temp\ai-bot-linux-amd64\install.sh
echo # Делаем исполняемым >> temp\ai-bot-linux-amd64\install.sh
echo chmod +x ./ai-bot >> temp\ai-bot-linux-amd64\install.sh
echo. >> temp\ai-bot-linux-amd64\install.sh
echo echo "✅ AI Bot успешно установлен!" >> temp\ai-bot-linux-amd64\install.sh
echo echo "" >> temp\ai-bot-linux-amd64\install.sh
echo echo "Быстрый старт:" >> temp\ai-bot-linux-amd64\install.sh
echo echo "  ./ai-bot --config  # Настройка" >> temp\ai-bot-linux-amd64\install.sh
echo echo "  ./ai-bot           # Запуск" >> temp\ai-bot-linux-amd64\install.sh

powershell -Command "Compress-Archive -Path 'temp\ai-bot-linux-amd64\*' -DestinationPath 'releases\ai-bot-linux-amd64.zip' -Force"
echo ✅ ai-bot-linux-amd64.zip

REM Linux 32-bit
echo Building Linux 32-bit...
mkdir temp\ai-bot-linux-386
set GOOS=linux
set GOARCH=386
go build -ldflags="-s -w" -o temp\ai-bot-linux-386\ai-bot .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-linux-386\ >nul 2>&1
copy USAGE.md temp\ai-bot-linux-386\ >nul 2>&1
copy .env.example temp\ai-bot-linux-386\ >nul 2>&1
copy temp\ai-bot-linux-amd64\install.sh temp\ai-bot-linux-386\ >nul

powershell -Command "Compress-Archive -Path 'temp\ai-bot-linux-386\*' -DestinationPath 'releases\ai-bot-linux-386.zip' -Force"
echo ✅ ai-bot-linux-386.zip

REM Linux ARM64
echo Building Linux ARM64...
mkdir temp\ai-bot-linux-arm64
set GOOS=linux
set GOARCH=arm64
go build -ldflags="-s -w" -o temp\ai-bot-linux-arm64\ai-bot .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-linux-arm64\ >nul 2>&1
copy USAGE.md temp\ai-bot-linux-arm64\ >nul 2>&1
copy .env.example temp\ai-bot-linux-arm64\ >nul 2>&1
copy temp\ai-bot-linux-amd64\install.sh temp\ai-bot-linux-arm64\ >nul

powershell -Command "Compress-Archive -Path 'temp\ai-bot-linux-arm64\*' -DestinationPath 'releases\ai-bot-linux-arm64.zip' -Force"
echo ✅ ai-bot-linux-arm64.zip

echo.
echo Building FreeBSD releases...

REM FreeBSD 64-bit
echo Building FreeBSD 64-bit...
mkdir temp\ai-bot-freebsd-amd64
set GOOS=freebsd
set GOARCH=amd64
go build -ldflags="-s -w" -o temp\ai-bot-freebsd-amd64\ai-bot .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-freebsd-amd64\ >nul 2>&1
copy USAGE.md temp\ai-bot-freebsd-amd64\ >nul 2>&1
copy .env.example temp\ai-bot-freebsd-amd64\ >nul 2>&1
copy temp\ai-bot-linux-amd64\install.sh temp\ai-bot-freebsd-amd64\ >nul

powershell -Command "Compress-Archive -Path 'temp\ai-bot-freebsd-amd64\*' -DestinationPath 'releases\ai-bot-freebsd-amd64.zip' -Force"
echo ✅ ai-bot-freebsd-amd64.zip

REM FreeBSD 32-bit
echo Building FreeBSD 32-bit...
mkdir temp\ai-bot-freebsd-386
set GOOS=freebsd
set GOARCH=386
go build -ldflags="-s -w" -o temp\ai-bot-freebsd-386\ai-bot .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-freebsd-386\ >nul 2>&1
copy USAGE.md temp\ai-bot-freebsd-386\ >nul 2>&1
copy .env.example temp\ai-bot-freebsd-386\ >nul 2>&1
copy temp\ai-bot-linux-amd64\install.sh temp\ai-bot-freebsd-386\ >nul

powershell -Command "Compress-Archive -Path 'temp\ai-bot-freebsd-386\*' -DestinationPath 'releases\ai-bot-freebsd-386.zip' -Force"
echo ✅ ai-bot-freebsd-386.zip

echo.
echo Building macOS releases...

REM macOS Intel
echo Building macOS Intel...
mkdir temp\ai-bot-darwin-amd64
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o temp\ai-bot-darwin-amd64\ai-bot .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-darwin-amd64\ >nul 2>&1
copy USAGE.md temp\ai-bot-darwin-amd64\ >nul 2>&1
copy .env.example temp\ai-bot-darwin-amd64\ >nul 2>&1
copy temp\ai-bot-linux-amd64\install.sh temp\ai-bot-darwin-amd64\ >nul

powershell -Command "Compress-Archive -Path 'temp\ai-bot-darwin-amd64\*' -DestinationPath 'releases\ai-bot-darwin-amd64.zip' -Force"
echo ✅ ai-bot-darwin-amd64.zip

REM macOS Apple Silicon
echo Building macOS Apple Silicon...
mkdir temp\ai-bot-darwin-arm64
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o temp\ai-bot-darwin-arm64\ai-bot .
if %errorlevel% neq 0 goto :error

copy README.md temp\ai-bot-darwin-arm64\ >nul 2>&1
copy USAGE.md temp\ai-bot-darwin-arm64\ >nul 2>&1
copy .env.example temp\ai-bot-darwin-arm64\ >nul 2>&1
copy temp\ai-bot-linux-amd64\install.sh temp\ai-bot-darwin-arm64\ >nul

powershell -Command "Compress-Archive -Path 'temp\ai-bot-darwin-arm64\*' -DestinationPath 'releases\ai-bot-darwin-arm64.zip' -Force"
echo ✅ ai-bot-darwin-arm64.zip

echo.
echo Generating checksums...
powershell -Command "Get-ChildItem 'releases\*.zip' | ForEach-Object { $hash = Get-FileHash $_.FullName -Algorithm SHA256; \"$($hash.Hash.ToLower())  $($_.Name)\" } | Set-Content -Path 'releases\checksums.txt' -Encoding UTF8"

echo.
echo Creating release README...
echo # AI Bot Release Files > releases\README.txt
echo. >> releases\README.txt
echo ## 📦 Скачать для вашей платформы: >> releases\README.txt
echo. >> releases\README.txt
echo ### Windows >> releases\README.txt
echo - **ai-bot-windows-amd64.zip** - Windows 64-bit (рекомендуется) >> releases\README.txt
echo - **ai-bot-windows-386.zip** - Windows 32-bit >> releases\README.txt
echo. >> releases\README.txt
echo ### Linux >> releases\README.txt
echo - **ai-bot-linux-amd64.zip** - Linux 64-bit (рекомендуется) >> releases\README.txt
echo - **ai-bot-linux-386.zip** - Linux 32-bit >> releases\README.txt
echo - **ai-bot-linux-arm64.zip** - Linux ARM64 (Raspberry Pi 4+) >> releases\README.txt
echo. >> releases\README.txt
echo ### FreeBSD >> releases\README.txt
echo - **ai-bot-freebsd-amd64.zip** - FreeBSD 64-bit >> releases\README.txt
echo - **ai-bot-freebsd-386.zip** - FreeBSD 32-bit >> releases\README.txt
echo. >> releases\README.txt
echo ### macOS >> releases\README.txt
echo - **ai-bot-darwin-amd64.zip** - macOS Intel >> releases\README.txt
echo - **ai-bot-darwin-arm64.zip** - macOS Apple Silicon (M1/M2/M3) >> releases\README.txt

REM Очищаем временные файлы
rmdir /S /Q temp 2>nul

echo.
echo 🎉 All releases completed successfully!
echo.
echo 📁 Built files:
dir releases /b
echo.
echo 🚀 Ready for GitHub release!
echo    Upload files from 'releases' folder
goto :end

:error
echo.
echo ❌ Build failed!
exit /b 1

:end