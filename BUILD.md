# 🔨 Сборка релизов AI Bot

Этот документ описывает как собрать релизы AI Bot для разных платформ.

## 🚀 Быстрый старт

### Полная сборка релизов (рекомендуется)

#### Windows (PowerShell)
```powershell
# Полная сборка с документацией и архивами
.\build-release.ps1

# Создание GitHub релиза (требует GitHub CLI)
.\create-github-release.ps1 -Version "v1.0.0"
```

#### Windows (Batch)
```cmd
# Полная сборка с документацией и архивами
build-release.bat
```

### Быстрая сборка (только исполняемые файлы)

#### Windows (PowerShell)
```powershell
# Сборка всех релизов с красивым выводом
.\build.ps1
```

#### Windows (Batch)
```cmd
# Сборка всех релизов
build.bat

# Быстрая сборка для текущей платформы
quick-build.bat
```

### Linux/macOS/FreeBSD (Make)
```bash
# Сборка всех релизов
make all

# Сборка только для Linux
make build-linux

# Сборка только для Windows
make build-windows

# Сборка только для FreeBSD
make build-freebsd

# Сборка только для macOS
make build-darwin

# Локальная сборка
make build

# Помощь
make help
```

## 📦 Поддерживаемые платформы

### Полные релизы (с документацией)

| Платформа | Архитектура | Архив | Содержимое |
|-----------|-------------|-------|------------|
| Windows | amd64 | `ai-bot-windows-amd64.zip` | ai-bot.exe + docs + install.bat |
| Windows | 386 | `ai-bot-windows-386.zip` | ai-bot.exe + docs + install.bat |
| Linux | amd64 | `ai-bot-linux-amd64.tar.gz` | ai-bot + docs + install.sh |
| Linux | 386 | `ai-bot-linux-386.tar.gz` | ai-bot + docs + install.sh |
| Linux | arm64 | `ai-bot-linux-arm64.tar.gz` | ai-bot + docs + install.sh |
| FreeBSD | amd64 | `ai-bot-freebsd-amd64.tar.gz` | ai-bot + docs + install.sh |
| FreeBSD | 386 | `ai-bot-freebsd-386.tar.gz` | ai-bot + docs + install.sh |
| macOS | amd64 | `ai-bot-darwin-amd64.tar.gz` | ai-bot + docs + install.sh |
| macOS | arm64 | `ai-bot-darwin-arm64.tar.gz` | ai-bot + docs + install.sh |

### Содержимое каждого архива:
- **ai-bot** (или ai-bot.exe) - исполняемый файл
- **README.md** - основная документация
- **USAGE.md** - руководство по использованию
- **.env.example** - пример конфигурации
- **install.sh** (Unix) или **install.bat** (Windows) - скрипт установки

## 🤖 Автоматическая сборка (GitHub Actions)

### Создание релиза по тегу
```bash
# Создаем тег и пушим
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions автоматически соберет и создаст релиз
```

### Ручной запуск сборки
1. Идите в GitHub → Actions
2. Выберите "Build and Release"
3. Нажмите "Run workflow"
4. Укажите версию (например, v1.0.1)
5. Нажмите "Run workflow"

## 🎯 Создание GitHub релиза локально

### Требования
- [GitHub CLI](https://cli.github.com/) установлен и настроен
- Права на создание релизов в репозитории

### Команды
```powershell
# 1. Собираем релизы
.\build-release.ps1

# 2. Создаем релиз на GitHub
.\create-github-release.ps1 -Version "v1.0.0"

# Дополнительные опции:
.\create-github-release.ps1 -Version "v1.0.0" -Title "Первый релиз" -Draft
.\create-github-release.ps1 -Version "v1.0.0-beta" -Prerelease
```

### Что происходит:
1. Создается тег в Git
2. Создается релиз на GitHub
3. Загружаются все архивы
4. Загружаются checksums.txt и README.txt
5. Генерируется красивое описание релиза

## 🛠️ Ручная сборка

### Требования
- Go 1.24 или новее
- Git (для версионирования)

### Команды сборки

#### Для Windows:
```bash
# 64-bit
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o releases/ai-bot-windows-amd64.exe .

# 32-bit
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o releases/ai-bot-windows-386.exe .
```

#### Для Linux:
```bash
# 64-bit
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o releases/ai-bot-linux-amd64 .

# 32-bit
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o releases/ai-bot-linux-386 .

# ARM64
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o releases/ai-bot-linux-arm64 .
```

#### Для FreeBSD:
```bash
# 64-bit
GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w" -o releases/ai-bot-freebsd-amd64 .

# 32-bit
GOOS=freebsd GOARCH=386 go build -ldflags "-s -w" -o releases/ai-bot-freebsd-386 .
```

#### Для macOS:
```bash
# Intel
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o releases/ai-bot-darwin-amd64 .

# Apple Silicon
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o releases/ai-bot-darwin-arm64 .
```

## 📋 Флаги сборки

- `-ldflags "-s -w"` - убирает отладочную информацию и таблицу символов (уменьшает размер)
- `-o <filename>` - указывает имя выходного файла

## 🔍 Проверка сборки

После сборки проверьте файлы:

```bash
# Размеры файлов
ls -lh releases/

# Проверка архитектуры (Linux/macOS)
file releases/ai-bot-linux-amd64

# Тест запуска
./releases/ai-bot-linux-amd64 --help
```

## 📦 Создание архивов

### Для релиза на GitHub:
```bash
cd releases

# Windows - ZIP архивы
zip ai-bot-windows-amd64.zip ai-bot-windows-amd64.exe
zip ai-bot-windows-386.zip ai-bot-windows-386.exe

# Unix системы - TAR.GZ архивы
tar -czf ai-bot-linux-amd64.tar.gz ai-bot-linux-amd64
tar -czf ai-bot-freebsd-amd64.tar.gz ai-bot-freebsd-amd64
tar -czf ai-bot-darwin-amd64.tar.gz ai-bot-darwin-amd64

# Контрольные суммы
sha256sum *.zip *.tar.gz > checksums.txt
```

## 🚨 Устранение проблем

### Ошибка "command not found: go"
Установите Go с официального сайта: https://golang.org/dl/

### Ошибка сборки для другой платформы
Убедитесь что установлены переменные окружения:
```bash
export GOOS=linux
export GOARCH=amd64
go build ...
```

### Большой размер файла
Используйте флаги оптимизации:
```bash
go build -ldflags "-s -w" -trimpath .
```

### Проблемы с зависимостями
```bash
go mod tidy
go mod download
```

## 📚 Дополнительные ресурсы

- [Go Cross Compilation](https://golang.org/doc/install/source#environment)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Build Modes](https://golang.org/cmd/go/#hdr-Build_modes)