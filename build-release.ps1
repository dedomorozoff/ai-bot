# AI Bot Release Builder - Полная версия с архивами
# Создает готовые к распространению архивы с документацией

Write-Host "🤖 Building AI Bot complete releases..." -ForegroundColor Cyan
Write-Host ""

# Создаем папки
if (!(Test-Path "releases")) { New-Item -ItemType Directory -Path "releases" | Out-Null }
if (!(Test-Path "temp")) { New-Item -ItemType Directory -Path "temp" | Out-Null }

# Очищаем старые релизы
Remove-Item "releases\*" -Force -ErrorAction SilentlyContinue
Remove-Item "temp\*" -Recurse -Force -ErrorAction SilentlyContinue

# Определяем платформы
$platforms = @(
    @{OS="windows"; Arch="amd64"; Ext=".exe"; Name="Windows 64-bit"; Archive="zip"},
    @{OS="windows"; Arch="386"; Ext=".exe"; Name="Windows 32-bit"; Archive="zip"},
    @{OS="linux"; Arch="amd64"; Ext=""; Name="Linux 64-bit"; Archive="tar.gz"},
    @{OS="linux"; Arch="386"; Ext=""; Name="Linux 32-bit"; Archive="tar.gz"},
    @{OS="linux"; Arch="arm64"; Ext=""; Name="Linux ARM64"; Archive="tar.gz"},
    @{OS="freebsd"; Arch="amd64"; Ext=""; Name="FreeBSD 64-bit"; Archive="tar.gz"},
    @{OS="freebsd"; Arch="386"; Ext=""; Name="FreeBSD 32-bit"; Archive="tar.gz"},
    @{OS="darwin"; Arch="amd64"; Ext=""; Name="macOS Intel"; Archive="tar.gz"},
    @{OS="darwin"; Arch="arm64"; Ext=""; Name="macOS Apple Silicon"; Archive="tar.gz"}
)

$success = 0
$total = $platforms.Count

foreach ($platform in $platforms) {
    $binaryName = "ai-bot$($platform.Ext)"
    $folderName = "ai-bot-$($platform.OS)-$($platform.Arch)"
    $tempDir = "temp\$folderName"
    
    Write-Host "🔨 Building $($platform.Name)..." -ForegroundColor Yellow
    
    # Создаем временную папку для этой платформы
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    # Собираем бинарник
    $env:GOOS = $platform.OS
    $env:GOARCH = $platform.Arch
    
    $buildArgs = @(
        "build",
        "-ldflags", "`"-s -w`"",
        "-o", "$tempDir\$binaryName",
        "."
    )
    
    $process = Start-Process -FilePath "go" -ArgumentList $buildArgs -Wait -PassThru -NoNewWindow
    
    if ($process.ExitCode -eq 0) {
        # Копируем дополнительные файлы
        Copy-Item "README.md" "$tempDir\" -ErrorAction SilentlyContinue
        Copy-Item "USAGE.md" "$tempDir\" -ErrorAction SilentlyContinue
        Copy-Item ".env.example" "$tempDir\" -ErrorAction SilentlyContinue
        
        # Создаем install скрипт для Unix систем
        if ($platform.OS -ne "windows") {
            $installScript = @"
#!/bin/bash
# AI Bot Installation Script

echo "🤖 AI Bot Installation"
echo "======================"
echo ""

# Делаем исполняемым
chmod +x ./ai-bot

# Проверяем установку
if ./ai-bot --help > /dev/null 2>&1; then
    echo "✅ AI Bot успешно установлен!"
    echo ""
    echo "Быстрый старт:"
    echo "  ./ai-bot --config  # Настройка"
    echo "  ./ai-bot           # Запуск"
    echo ""
else
    echo "❌ Ошибка установки"
    exit 1
fi
"@
            Set-Content -Path "$tempDir\install.sh" -Value $installScript -Encoding UTF8
        } else {
            # Создаем install.bat для Windows
            $installBat = @"
@echo off
echo 🤖 AI Bot Installation
echo ======================
echo.

echo ✅ AI Bot готов к использованию!
echo.
echo Быстрый старт:
echo   ai-bot.exe --config  # Настройка
echo   ai-bot.exe           # Запуск
echo.
pause
"@
            Set-Content -Path "$tempDir\install.bat" -Value $installBat -Encoding UTF8
        }
        
        # Создаем архив
        if ($platform.Archive -eq "zip") {
            $archiveName = "$folderName.zip"
            Compress-Archive -Path "$tempDir\*" -DestinationPath "releases\$archiveName" -Force
        } else {
            $archiveName = "$folderName.tar.gz"
            # Используем tar для создания tar.gz (если доступен)
            if (Get-Command tar -ErrorAction SilentlyContinue) {
                Set-Location $tempDir
                tar -czf "..\..\releases\$archiveName" *
                Set-Location ..\..
            } else {
                # Fallback к zip если tar недоступен
                $archiveName = "$folderName.zip"
                Compress-Archive -Path "$tempDir\*" -DestinationPath "releases\$archiveName" -Force
            }
        }
        
        Write-Host "   ✅ Success: $archiveName" -ForegroundColor Green
        $success++
    } else {
        Write-Host "   ❌ Failed: $folderName" -ForegroundColor Red
    }
}

# Создаем checksums
Write-Host ""
Write-Host "🔐 Generating checksums..." -ForegroundColor Yellow

$checksums = @()
Get-ChildItem "releases\*" | ForEach-Object {
    $hash = Get-FileHash $_.FullName -Algorithm SHA256
    $checksums += "$($hash.Hash.ToLower())  $($_.Name)"
}

Set-Content -Path "releases\checksums.txt" -Value $checksums -Encoding UTF8

# Создаем README для релиза
$releaseReadme = @"
# AI Bot Release Files

## 📦 Скачать для вашей платформы:

### Windows
- **ai-bot-windows-amd64.zip** - Windows 64-bit (рекомендуется)
- **ai-bot-windows-386.zip** - Windows 32-bit

### Linux
- **ai-bot-linux-amd64.tar.gz** - Linux 64-bit (рекомендуется)
- **ai-bot-linux-386.tar.gz** - Linux 32-bit
- **ai-bot-linux-arm64.tar.gz** - Linux ARM64 (Raspberry Pi 4+)

### FreeBSD
- **ai-bot-freebsd-amd64.tar.gz** - FreeBSD 64-bit
- **ai-bot-freebsd-386.tar.gz** - FreeBSD 32-bit

### macOS
- **ai-bot-darwin-amd64.tar.gz** - macOS Intel
- **ai-bot-darwin-arm64.tar.gz** - macOS Apple Silicon (M1/M2/M3)

## 🚀 Установка и запуск:

### Windows:
1. Скачайте ai-bot-windows-amd64.zip
2. Распакуйте архив
3. Запустите install.bat или:
   ```cmd
   ai-bot.exe --config
   ai-bot.exe
   ```

### Linux/macOS/FreeBSD:
1. Скачайте архив для вашей платформы
2. Распакуйте: `tar -xzf ai-bot-*.tar.gz`
3. Запустите: `./install.sh` или:
   ```bash
   chmod +x ai-bot
   ./ai-bot --config
   ./ai-bot
   ```

## 🔐 Проверка целостности:

Проверьте контрольные суммы:
```bash
sha256sum -c checksums.txt
```

## 📚 Документация:

- **README.md** - Основная документация
- **USAGE.md** - Руководство по использованию
- **.env.example** - Пример конфигурации

## 🆘 Поддержка:

Если возникли проблемы, создайте issue в GitHub репозитории.
"@

Set-Content -Path "releases\README.txt" -Value $releaseReadme -Encoding UTF8

# Очищаем временные файлы
Remove-Item "temp" -Recurse -Force -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "📊 Build Summary:" -ForegroundColor Cyan
Write-Host "   Successful: $success/$total" -ForegroundColor Green

if ($success -eq $total) {
    Write-Host "🎉 All releases completed successfully!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Some builds failed!" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "📁 Release files:" -ForegroundColor Cyan
Get-ChildItem "releases" | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "   $($_.Name) ($size MB)" -ForegroundColor White
}

Write-Host ""
Write-Host "🚀 Ready for GitHub release!" -ForegroundColor Green
Write-Host "   Upload files from 'releases' folder" -ForegroundColor Gray

# Сброс переменных окружения
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue