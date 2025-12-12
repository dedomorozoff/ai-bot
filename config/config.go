package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config структура конфигурации
type Config struct {
	Host           string
	Port           string
	OpenRouterKey  string
	OpenRouterModel string
	OpenAIKey      string
	OpenAIModel    string
	MaxTokens      int
	Temperature    float64
	Timeout        int
	SystemPrompt   string
}

// Load загружает конфигурацию из .env файла и переменных окружения
func Load() (*Config, error) {
	// Загружаем .env файл если он существует
	_ = godotenv.Load()

	cfg := &Config{
		Host:           getEnv("HOST", "0.0.0.0"),
		Port:           getEnv("PORT", "8080"),
		OpenRouterKey:  getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterModel: getEnv("OPENROUTER_MODEL", "anthropic/claude-3.5-sonnet"),
		OpenAIKey:      getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:    getEnv("OPENAI_MODEL", "gpt-4o"),
		MaxTokens:      getEnvInt("MAX_TOKENS", 4000),
		Temperature:    getEnvFloat("TEMPERATURE", 0.3),
		Timeout:        getEnvInt("TIMEOUT", 30),
		SystemPrompt:   getEnv("SYSTEM_PROMPT", "Ты полезный AI ассистент. Отвечай кратко и по делу на русском языке."),
	}

	return cfg, nil
}

// Save сохраняет конфигурацию в .env файл
func Save(cfg *Config) error {
	envFile := ".env"
	
	// Читаем существующий файл если он есть
	existing := make(map[string]string)
	if file, err := os.Open(envFile); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				existing[parts[0]] = parts[1]
			}
		}
		file.Close()
	}

	// Обновляем значения
	existing["HOST"] = cfg.Host
	existing["PORT"] = cfg.Port
	existing["OPENROUTER_API_KEY"] = cfg.OpenRouterKey
	existing["OPENROUTER_MODEL"] = cfg.OpenRouterModel
	existing["OPENAI_API_KEY"] = cfg.OpenAIKey
	existing["OPENAI_MODEL"] = cfg.OpenAIModel
	existing["MAX_TOKENS"] = strconv.Itoa(cfg.MaxTokens)
	existing["TEMPERATURE"] = fmt.Sprintf("%.2f", cfg.Temperature)
	existing["TIMEOUT"] = strconv.Itoa(cfg.Timeout)
	existing["SYSTEM_PROMPT"] = cfg.SystemPrompt

	// Записываем обратно
	file, err := os.Create(envFile)
	if err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	
	// Записываем в определенном порядке
	keys := []string{
		"HOST", "PORT",
		"OPENROUTER_API_KEY", "OPENROUTER_MODEL",
		"OPENAI_API_KEY", "OPENAI_MODEL",
		"MAX_TOKENS", "TEMPERATURE", "TIMEOUT",
		"SYSTEM_PROMPT",
	}

	for _, key := range keys {
		if value, ok := existing[key]; ok && value != "" {
			fmt.Fprintf(writer, "%s=%s\n", key, value)
		}
	}

	// Записываем остальные ключи
	for key, value := range existing {
		found := false
		for _, k := range keys {
			if k == key {
				found = true
				break
			}
		}
		if !found && value != "" {
			fmt.Fprintf(writer, "%s=%s\n", key, value)
		}
	}

	return writer.Flush()
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt получает int значение переменной окружения
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvFloat получает float64 значение переменной окружения
func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

