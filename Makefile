# AI Bot Makefile
# Сборка релизов для разных платформ

.PHONY: all clean build-windows build-linux build-freebsd build-darwin releases

# Версия приложения (можно передать через VERSION=x.x.x)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Флаги сборки
LDFLAGS = -s -w -X main.version=$(VERSION)
BUILD_FLAGS = -ldflags "$(LDFLAGS)"

# Создание папки для релизов
releases:
	@mkdir -p releases
	@rm -f releases/*

# Сборка для Windows
build-windows: releases
	@echo "🔨 Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o releases/ai-bot-windows-amd64.exe .
	@GOOS=windows GOARCH=386 go build $(BUILD_FLAGS) -o releases/ai-bot-windows-386.exe .
	@echo "✅ Windows builds completed"

# Сборка для Linux
build-linux: releases
	@echo "🔨 Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o releases/ai-bot-linux-amd64 .
	@GOOS=linux GOARCH=386 go build $(BUILD_FLAGS) -o releases/ai-bot-linux-386 .
	@GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o releases/ai-bot-linux-arm64 .
	@echo "✅ Linux builds completed"

# Сборка для FreeBSD
build-freebsd: releases
	@echo "🔨 Building for FreeBSD..."
	@GOOS=freebsd GOARCH=amd64 go build $(BUILD_FLAGS) -o releases/ai-bot-freebsd-amd64 .
	@GOOS=freebsd GOARCH=386 go build $(BUILD_FLAGS) -o releases/ai-bot-freebsd-386 .
	@echo "✅ FreeBSD builds completed"

# Сборка для macOS
build-darwin: releases
	@echo "🔨 Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o releases/ai-bot-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o releases/ai-bot-darwin-arm64 .
	@echo "✅ macOS builds completed"

# Сборка всех релизов
all: build-windows build-linux build-freebsd build-darwin
	@echo ""
	@echo "🎉 All builds completed successfully!"
	@echo ""
	@echo "📁 Built files:"
	@ls -lh releases/
	@echo ""
	@echo "🚀 Ready for GitHub release!"

# Очистка
clean:
	@rm -rf releases/
	@echo "🧹 Cleaned up releases directory"

# Локальная сборка для текущей платформы
build:
	@echo "🔨 Building for current platform..."
	@go build $(BUILD_FLAGS) -o ai-bot .
	@echo "✅ Local build completed"

# Тестирование
test:
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅ Tests completed"

# Форматирование кода
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Проверка кода
vet:
	@echo "🔍 Vetting code..."
	@go vet ./...
	@echo "✅ Code vetted"

# Полная проверка перед релизом
check: fmt vet test
	@echo "✅ All checks passed"

# Помощь
help:
	@echo "AI Bot Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all           - Build all releases"
	@echo "  build         - Build for current platform"
	@echo "  build-windows - Build Windows releases"
	@echo "  build-linux   - Build Linux releases"
	@echo "  build-freebsd - Build FreeBSD releases"
	@echo "  build-darwin  - Build macOS releases"
	@echo "  test          - Run tests"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  check         - Run all checks"
	@echo "  clean         - Clean releases"
	@echo "  help          - Show this help"
	@echo ""
	@echo "Usage examples:"
	@echo "  make all                    # Build all releases"
	@echo "  make all VERSION=1.0.0      # Build with specific version"
	@echo "  make build-linux            # Build only Linux releases"