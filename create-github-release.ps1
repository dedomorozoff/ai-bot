# Скрипт для создания GitHub релиза
# Требует установленного GitHub CLI (gh)

param(
    [Parameter(Mandatory=$true)]
    [string]$Version,
    
    [Parameter(Mandatory=$false)]
    [string]$Title = "",
    
    [Parameter(Mandatory=$false)]
    [switch]$Draft = $false,
    
    [Parameter(Mandatory=$false)]
    [switch]$Prerelease = $false
)

Write-Host "🚀 Creating GitHub release $Version..." -ForegroundColor Cyan
Write-Host ""

# Проверяем наличие GitHub CLI
if (!(Get-Command gh -ErrorAction SilentlyContinue)) {
    Write-Host "❌ GitHub CLI (gh) не найден!" -ForegroundColor Red
    Write-Host "   Установите с https://cli.github.com/" -ForegroundColor Yellow
    exit 1
}

# Проверяем наличие файлов релиза
if (!(Test-Path "releases")) {
    Write-Host "❌ Папка releases не найдена!" -ForegroundColor Red
    Write-Host "   Запустите сначала build-release.ps1" -ForegroundColor Yellow
    exit 1
}

$releaseFiles = Get-ChildItem "releases\*" -File
if ($releaseFiles.Count -eq 0) {
    Write-Host "❌ Файлы релиза не найдены!" -ForegroundColor Red
    Write-Host "   Запустите сначала build-release.ps1" -ForegroundColor Yellow
    exit 1
}

# Устанавливаем заголовок по умолчанию
if ($Title -eq "") {
    $Title = "AI Bot $Version"
}

# Создаем описание релиза
$releaseNotes = @"
## AI Bot $Version

### 📦 Скачать релиз:

**Windows:**
- [ai-bot-windows-amd64.zip](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-windows-amd64.zip) - Windows 64-bit (рекомендуется)
- [ai-bot-windows-386.zip](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-windows-386.zip) - Windows 32-bit

**Linux:**
- [ai-bot-linux-amd64.tar.gz](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-linux-amd64.tar.gz) - Linux 64-bit (рекомендуется)
- [ai-bot-linux-386.tar.gz](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-linux-386.tar.gz) - Linux 32-bit
- [ai-bot-linux-arm64.tar.gz](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-linux-arm64.tar.gz) - Linux ARM64 (Raspberry Pi 4+)

**FreeBSD:**
- [ai-bot-freebsd-amd64.tar.gz](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-freebsd-amd64.tar.gz) - FreeBSD 64-bit
- [ai-bot-freebsd-386.tar.gz](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-freebsd-386.tar.gz) - FreeBSD 32-bit

**macOS:**
- [ai-bot-darwin-amd64.tar.gz](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-darwin-amd64.tar.gz) - macOS Intel
- [ai-bot-darwin-arm64.tar.gz](https://github.com/`${{ github.repository }}/releases/download/$Version/ai-bot-darwin-arm64.tar.gz) - macOS Apple Silicon (M1/M2/M3)

### 🚀 Быстрый старт:

#### Windows:
1. Скачайте ai-bot-windows-amd64.zip
2. Распакуйте архив
3. Запустите install.bat или:
   ```cmd
   ai-bot.exe --config
   ai-bot.exe
   ```

#### Linux/macOS/FreeBSD:
1. Скачайте архив для вашей платформы
2. Распакуйте: ``tar -xzf ai-bot-*.tar.gz``
3. Запустите: ``./install.sh`` или:
   ```bash
   chmod +x ai-bot
   ./ai-bot --config
   ./ai-bot
   ```

### 📋 Что включено в каждый архив:

- **ai-bot** - исполняемый файл
- **README.md** - основная документация
- **USAGE.md** - руководство по использованию
- **.env.example** - пример конфигурации
- **install** скрипт для быстрой установки

### 🔐 Контрольные суммы:

Проверьте целостность файлов с помощью [checksums.txt](https://github.com/`${{ github.repository }}/releases/download/$Version/checksums.txt)

```bash
sha256sum -c checksums.txt
```

### 🆘 Поддержка:

Если возникли проблемы, создайте [issue](https://github.com/`${{ github.repository }}/issues) в репозитории.
"@

# Сохраняем описание в файл
Set-Content -Path "release-notes.md" -Value $releaseNotes -Encoding UTF8

# Формируем команду gh
$ghArgs = @(
    "release", "create", $Version,
    "--title", $Title,
    "--notes-file", "release-notes.md"
)

if ($Draft) {
    $ghArgs += "--draft"
}

if ($Prerelease) {
    $ghArgs += "--prerelease"
}

# Добавляем все файлы из папки releases
$releaseFiles | ForEach-Object {
    $ghArgs += $_.FullName
}

Write-Host "📝 Создаем релиз с параметрами:" -ForegroundColor Yellow
Write-Host "   Версия: $Version" -ForegroundColor White
Write-Host "   Заголовок: $Title" -ForegroundColor White
Write-Host "   Черновик: $Draft" -ForegroundColor White
Write-Host "   Пререлиз: $Prerelease" -ForegroundColor White
Write-Host "   Файлов: $($releaseFiles.Count)" -ForegroundColor White
Write-Host ""

# Создаем релиз
Write-Host "🚀 Создаем релиз..." -ForegroundColor Cyan
& gh @ghArgs

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "✅ Релиз $Version успешно создан!" -ForegroundColor Green
    Write-Host ""
    Write-Host "🔗 Ссылки:" -ForegroundColor Cyan
    Write-Host "   Релиз: https://github.com/$(gh repo view --json owner,name -q '.owner.login + "/" + .name")/releases/tag/$Version" -ForegroundColor White
    Write-Host "   Все релизы: https://github.com/$(gh repo view --json owner,name -q '.owner.login + "/" + .name")/releases" -ForegroundColor White
} else {
    Write-Host ""
    Write-Host "❌ Ошибка создания релиза!" -ForegroundColor Red
    Write-Host "   Проверьте права доступа и подключение к GitHub" -ForegroundColor Yellow
}

# Очищаем временный файл
Remove-Item "release-notes.md" -ErrorAction SilentlyContinue