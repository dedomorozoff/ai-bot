# AI Bot Release Builder
# Собирает релизы для Windows, Linux, FreeBSD и macOS

Write-Host "🤖 Building AI Bot releases for multiple platforms..." -ForegroundColor Cyan
Write-Host ""

# Создаем папку для релизов
if (!(Test-Path "releases")) {
    New-Item -ItemType Directory -Path "releases" | Out-Null
}

# Очищаем старые релизы
Remove-Item "releases\*" -Force -ErrorAction SilentlyContinue

# Определяем платформы для сборки
$platforms = @(
    @{OS="windows"; Arch="amd64"; Ext=".exe"; Name="Windows 64-bit"},
    @{OS="windows"; Arch="386"; Ext=".exe"; Name="Windows 32-bit"},
    @{OS="linux"; Arch="amd64"; Ext=""; Name="Linux 64-bit"},
    @{OS="linux"; Arch="386"; Ext=""; Name="Linux 32-bit"},
    @{OS="linux"; Arch="arm64"; Ext=""; Name="Linux ARM64"},
    @{OS="freebsd"; Arch="amd64"; Ext=""; Name="FreeBSD 64-bit"},
    @{OS="freebsd"; Arch="386"; Ext=""; Name="FreeBSD 32-bit"},
    @{OS="darwin"; Arch="amd64"; Ext=""; Name="macOS Intel"},
    @{OS="darwin"; Arch="arm64"; Ext=""; Name="macOS Apple Silicon"}
)

$success = 0
$total = $platforms.Count

foreach ($platform in $platforms) {
    $outputName = "ai-bot-$($platform.OS)-$($platform.Arch)$($platform.Ext)"
    
    Write-Host "🔨 Building $($platform.Name)..." -ForegroundColor Yellow
    
    $env:GOOS = $platform.OS
    $env:GOARCH = $platform.Arch
    
    $buildArgs = @(
        "build",
        "-ldflags", "`"-s -w`"",
        "-o", "releases\$outputName",
        "."
    )
    
    $process = Start-Process -FilePath "go" -ArgumentList $buildArgs -Wait -PassThru -NoNewWindow
    
    if ($process.ExitCode -eq 0) {
        Write-Host "   ✅ Success: $outputName" -ForegroundColor Green
        $success++
    } else {
        Write-Host "   ❌ Failed: $outputName" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "📊 Build Summary:" -ForegroundColor Cyan
Write-Host "   Successful: $success/$total" -ForegroundColor Green

if ($success -eq $total) {
    Write-Host "🎉 All builds completed successfully!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Some builds failed!" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "📁 Built files:" -ForegroundColor Cyan
Get-ChildItem "releases" | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "   $($_.Name) ($size MB)" -ForegroundColor White
}

Write-Host ""
Write-Host "🚀 Ready for GitHub release!" -ForegroundColor Green

# Сброс переменных окружения
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue