package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config конфигурация AI клиента
type Config struct {
	OpenRouterAPIKey string
	OpenRouterModel  string
	OpenRouterURL    string
	OpenAIAPIKey     string
	OpenAIModel      string
	MaxTokens        int
	Temperature      float32
	RequestTimeout   int
}

// ChatMessage представляет сообщение в чате
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Client клиент для работы с AI
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient создает новый AI клиент
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
	}
}

// Chat отправляет запрос в чат с AI
func (c *Client) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	// Пробуем OpenRouter сначала
	if c.config.OpenRouterAPIKey != "" {
		response, err := c.chatOpenRouter(ctx, messages)
		if err == nil {
			return response, nil
		}
		// Логируем ошибку, но продолжаем с fallback
		fmt.Printf("OpenRouter failed, falling back to OpenAI: %v\n", err)
	}

	// Fallback на OpenAI
	if c.config.OpenAIAPIKey != "" {
		return c.chatOpenAI(ctx, messages)
	}

	return "", fmt.Errorf("no AI provider configured")
}

// chatOpenRouter отправляет запрос в OpenRouter API
func (c *Client) chatOpenRouter(ctx context.Context, messages []ChatMessage) (string, error) {
	openAIMessages := make([]openAIMessage, len(messages))
	for i, msg := range messages {
		openAIMessages[i] = openAIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	request := openAIRequest{
		Model:       c.config.OpenRouterModel,
		Messages:    openAIMessages,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := c.config.OpenRouterURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.OpenRouterAPIKey)
	req.Header.Set("HTTP-Referer", "https://ai-bot.local")
	req.Header.Set("X-Title", "AI Bot")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(responseBody, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenRouter API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenRouter")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// chatOpenAI отправляет запрос в OpenAI API
func (c *Client) chatOpenAI(ctx context.Context, messages []ChatMessage) (string, error) {
	openAIMessages := make([]openAIMessage, len(messages))
	for i, msg := range messages {
		openAIMessages[i] = openAIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	request := openAIRequest{
		Model:       c.config.OpenAIModel,
		Messages:    openAIMessages,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.OpenAIAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(responseBody, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// IsConfigured проверяет, настроен ли AI клиент
func (c *Client) IsConfigured() bool {
	return c.config.OpenRouterAPIKey != "" || c.config.OpenAIAPIKey != ""
}

// GetProvider возвращает текущего провайдера AI
func (c *Client) GetProvider() string {
	if c.config.OpenRouterAPIKey != "" {
		return fmt.Sprintf("OpenRouter (%s)", c.config.OpenRouterModel)
	}
	if c.config.OpenAIAPIKey != "" {
		return fmt.Sprintf("OpenAI (%s)", c.config.OpenAIModel)
	}
	return "Not configured"
}

// ModelInfo информация о модели
type ModelInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ContextLength int    `json:"context_length"`
	Pricing       *struct {
		Prompt     interface{} `json:"prompt"`
		Completion interface{} `json:"completion"`
	} `json:"pricing,omitempty"`
}

// OpenRouterModelsResponse ответ от OpenRouter API с моделями
type OpenRouterModelsResponse struct {
	Data []ModelInfo `json:"data"`
}

// GetModels получает список доступных моделей от OpenRouter API
func (c *Client) GetModels(ctx context.Context) ([]ModelInfo, error) {
	if c.config.OpenRouterAPIKey == "" {
		return nil, fmt.Errorf("OpenRouter API key not configured")
	}

	apiURL := c.config.OpenRouterURL + "/models"
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.OpenRouterAPIKey)
	req.Header.Set("HTTP-Referer", "https://ai-bot.local")
	req.Header.Set("X-Title", "AI Bot")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var modelsResp OpenRouterModelsResponse
	if err := json.Unmarshal(responseBody, &modelsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return modelsResp.Data, nil
}

// Внутренние структуры для API
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float32         `json:"temperature"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}
