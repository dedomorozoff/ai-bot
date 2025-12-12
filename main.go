package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"ai-bot/ai"
	"ai-bot/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	host            = flag.String("host", "", "Host to bind server (overrides .env)")
	port            = flag.String("port", "", "Port to bind server (overrides .env)")
	openRouterKey   = flag.String("openrouter-key", "", "OpenRouter API key (overrides .env)")
	openRouterModel = flag.String("openrouter-model", "", "OpenRouter model (overrides .env)")
	openAIKey       = flag.String("openai-key", "", "OpenAI API key (overrides .env)")
	openAIModel     = flag.String("openai-model", "", "OpenAI model (overrides .env)")
	maxTokens       = flag.Int("max-tokens", 0, "Maximum tokens (overrides .env)")
	temperature     = flag.Float64("temperature", -1, "Temperature 0-1 (overrides .env)")
	timeout         = flag.Int("timeout", 0, "Request timeout in seconds (overrides .env)")
	configCmd       = flag.Bool("config", false, "Run configuration wizard")
	demoOnly        = flag.Bool("demo", false, "Show demo page on main route (/)")
)

func main() {
	flag.Parse()

	// Если указан флаг config, запускаем конфигуратор
	if *configCmd {
		runConfig()
		return
	}

	// Загружаем конфигурацию из .env
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Предупреждение: не удалось загрузить .env файл: %v", err)
		cfg = &config.Config{
			Host:            "0.0.0.0",
			Port:            "8080",
			OpenRouterModel: "anthropic/claude-3.5-sonnet",
			OpenAIModel:     "gpt-4o",
			MaxTokens:       4000,
			Temperature:     0.3,
			Timeout:         30,
		}
	}

	// Переопределяем значения из командной строки
	if *host != "" {
		cfg.Host = *host
	}
	if *port != "" {
		cfg.Port = *port
	}
	if *openRouterKey != "" {
		cfg.OpenRouterKey = *openRouterKey
	}
	if *openRouterModel != "" {
		cfg.OpenRouterModel = *openRouterModel
	}
	if *openAIKey != "" {
		cfg.OpenAIKey = *openAIKey
	}
	if *openAIModel != "" {
		cfg.OpenAIModel = *openAIModel
	}
	if *maxTokens > 0 {
		cfg.MaxTokens = *maxTokens
	}
	if *temperature >= 0 {
		cfg.Temperature = *temperature
	}
	if *timeout > 0 {
		cfg.Timeout = *timeout
	}

	// Создаем конфигурацию AI
	aiConfig := &ai.Config{
		OpenRouterAPIKey: cfg.OpenRouterKey,
		OpenRouterModel:  cfg.OpenRouterModel,
		OpenRouterURL:    "https://openrouter.ai/api/v1",
		OpenAIAPIKey:     cfg.OpenAIKey,
		OpenAIModel:      cfg.OpenAIModel,
		MaxTokens:        cfg.MaxTokens,
		Temperature:      float32(cfg.Temperature),
		RequestTimeout:   cfg.Timeout,
	}

	// Проверяем наличие хотя бы одного API ключа
	if aiConfig.OpenRouterAPIKey == "" && aiConfig.OpenAIAPIKey == "" {
		log.Fatal("Ошибка: необходимо указать хотя бы один API ключ")
		log.Fatal("Используйте: ./ai-bot.exe --config")
		log.Fatal("Или укажите ключ в .env файле или через аргументы командной строки")
	}

	// Создаем AI клиент
	client := ai.NewClient(aiConfig)

	// Настраиваем маршруты
	if *demoOnly {
		// Если указан флаг --demo, показываем демо страницу на главной
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			serveDemoPage(w, r)
		})
	} else {
		// Обычная главная страница с чатом
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			serveHTML(w, r)
		})

		// Демо страница на отдельном маршруте
		http.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
			serveDemoPage(w, r)
		})

		// Демо с кастомными цветами
		http.HandleFunc("/demo-custom", func(w http.ResponseWriter, r *http.Request) {
			serveCustomDemoPage(w, r)
		})

		// Демо с кастомным промптом
		http.HandleFunc("/demo-prompt", func(w http.ResponseWriter, r *http.Request) {
			servePromptDemoPage(w, r)
		})
	}

	http.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		handleChat(w, r, client, aiConfig)
	})

	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		handleStatus(w, r, client)
	})

	// Статические файлы (JS)
	http.HandleFunc("/static/ai-bot.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/ai-bot.js")
	})

	// Встроенный чат - один тег script
	http.HandleFunc("/chat.js", func(w http.ResponseWriter, r *http.Request) {
		serveEmbeddedChatNew(w, r)
	})

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("AI Bot сервер запущен на http://%s", addr)
	if *demoOnly {
		log.Printf("  Режим: Демо страница (--demo)")
	}
	log.Printf("Конфигурация:")
	log.Printf("  OpenRouter: %s", maskKey(cfg.OpenRouterKey))
	log.Printf("  OpenAI: %s", maskKey(cfg.OpenAIKey))
	if cfg.OpenRouterKey != "" {
		log.Printf("  Модель: %s", cfg.OpenRouterModel)
	} else if cfg.OpenAIKey != "" {
		log.Printf("  Модель: %s", cfg.OpenAIModel)
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

func maskKey(key string) string {
	if key == "" {
		return "не указан"
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Bot</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            width: 100%;
            max-width: 800px;
            height: 600px;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            text-align: center;
        }
        .header h1 {
            font-size: 24px;
            margin-bottom: 5px;
        }
        .header p {
            opacity: 0.9;
            font-size: 14px;
        }
        .chat-container {
            flex: 1;
            overflow-y: auto;
            padding: 20px;
            background: #f5f7fa;
        }
        .message {
            margin-bottom: 15px;
            display: flex;
            animation: fadeIn 0.3s;
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .message.user {
            justify-content: flex-end;
        }
        .message-content {
            max-width: 70%;
            padding: 12px 16px;
            border-radius: 18px;
            word-wrap: break-word;
        }
        .message.user .message-content {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        .message.bot .message-content {
            background: white;
            color: #333;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        .input-container {
            padding: 20px;
            background: white;
            border-top: 1px solid #e0e0e0;
            display: flex;
            gap: 10px;
        }
        #messageInput {
            flex: 1;
            padding: 12px 16px;
            border: 2px solid #e0e0e0;
            border-radius: 25px;
            font-size: 14px;
            outline: none;
            transition: border-color 0.3s;
        }
        #messageInput:focus {
            border-color: #667eea;
        }
        #sendButton {
            padding: 12px 24px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 25px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 600;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        #sendButton:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        #sendButton:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            transform: none;
        }
        .status {
            padding: 10px 20px;
            background: #f0f0f0;
            text-align: center;
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🤖 AI Bot</h1>
            <p>Чат с AI ассистентом</p>
            <p style="margin-top: 10px; font-size: 14px;">
                <a href="/demo" style="color: white; text-decoration: underline;">Демо страница</a>
            </p>
        </div>
        <div class="status" id="status">Подключение...</div>
        <div class="chat-container" id="chatContainer"></div>
        <div class="input-container">
            <input type="text" id="messageInput" placeholder="Введите сообщение..." autocomplete="off">
            <button id="sendButton">Отправить</button>
        </div>
    </div>
    <script src="/static/ai-bot.js"></script>
</body>
</html>`

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveDemoPage(w http.ResponseWriter, r *http.Request) {
	// Получаем базовый URL для скрипта
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	tmpl := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Bot - Демо страница</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .header {
            text-align: center;
            color: white;
            margin-bottom: 50px;
        }
        .header h1 {
            font-size: 48px;
            margin-bottom: 10px;
            text-shadow: 0 2px 10px rgba(0,0,0,0.2);
        }
        .header p {
            font-size: 20px;
            opacity: 0.9;
        }
        .demo-section {
            background: white;
            border-radius: 20px;
            padding: 40px;
            margin-bottom: 30px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
        }
        .demo-section h2 {
            color: #667eea;
            margin-bottom: 20px;
            font-size: 28px;
        }
        .demo-section p {
            color: #666;
            line-height: 1.6;
            margin-bottom: 15px;
            font-size: 16px;
        }
        .code-block {
            background: #f5f7fa;
            border: 1px solid #e0e0e0;
            border-radius: 10px;
            padding: 20px;
            margin: 20px 0;
            overflow-x: auto;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            color: #333;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            text-decoration: none;
            border-radius: 25px;
            font-weight: 600;
            margin: 10px 10px 10px 0;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        .button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        .footer {
            text-align: center;
            color: white;
            margin-top: 50px;
            opacity: 0.8;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🤖 AI Bot</h1>
            <p>Демо страница с примером интеграции</p>
        </div>

        <div class="demo-section">
            <h2>📋 О проекте</h2>
            <p>AI Bot - это простое приложение для чата с AI через OpenRouter или OpenAI API.</p>
            <p>Вы можете легко интегрировать его на любую HTML страницу всего одним тегом script.</p>
        </div>

        <div class="demo-section">
            <h2>🚀 Быстрый старт</h2>
            <p><strong>1. Настройка:</strong></p>
            <div class="code-block">
./ai-bot.exe --config
            </div>
            
            <p><strong>2. Запуск:</strong></p>
            <div class="code-block">
./ai-bot.exe
            </div>
            
            <p><strong>3. Интеграция в HTML:</strong></p>
            <div class="code-block">
&lt;script src="http://localhost:8080/chat.js"&gt;&lt;/script&gt;
            </div>
        </div>

        <div class="demo-section">
            <h2>💻 Базовая интеграция</h2>
            <p>Простейший способ подключения:</p>
            <div class="code-block">
&lt;script src="` + baseURL + `/chat.js"&gt;&lt;/script&gt;
            </div>
        </div>

        <div class="demo-section">
            <h2>🎨 Кастомизация цветов</h2>
            <p>Настройте цвета под ваш дизайн:</p>
            <div class="code-block">
&lt;script src="` + baseURL + `/chat.js"
        data-primary-color="#ff6b6b"
        data-secondary-color="#4ecdc4" 
        data-accent-color="#45b7d1"&gt;&lt;/script&gt;
            </div>
        </div>

        <div class="demo-section">
            <h2>🤖 Кастомный системный промпт</h2>
            <p>Настройте поведение AI под ваши задачи:</p>
            <div class="code-block">
&lt;script src="` + baseURL + `/chat.js"
        data-system-prompt="Ты дружелюбный помощник программиста. Объясняй код простыми словами и добавляй эмодзи."&gt;&lt;/script&gt;
            </div>
        </div>

        <div class="demo-section">
            <h2>⚙️ Продвинутая кастомизация</h2>
            <p>Комбинируйте все возможности:</p>
            <div class="code-block">
&lt;script src="` + baseURL + `/chat.js"
        data-primary-color="#2c3e50"
        data-secondary-color="#34495e"
        data-system-prompt="Ты профессиональный консультант. Давай четкие ответы."
        data-custom-css=".ai-chat-toggle{border:3px solid gold;}"&gt;&lt;/script&gt;
            </div>
        </div>

        <div class="demo-section">
            <h2>📋 Полный пример</h2>
            <div class="code-block">
&lt;!DOCTYPE html&gt;
&lt;html&gt;
&lt;head&gt;
    &lt;title&gt;Моя страница&lt;/title&gt;
&lt;/head&gt;
&lt;body&gt;
    &lt;h1&gt;Добро пожаловать!&lt;/h1&gt;
    &lt;p&gt;Любой контент вашей страницы...&lt;/p&gt;
    
    &lt;!-- AI Bot с кастомными цветами --&gt;
    &lt;script src="` + baseURL + `/chat.js"
            data-primary-color="#e74c3c"
            data-secondary-color="#c0392b"&gt;&lt;/script&gt;
&lt;/body&gt;
&lt;/html&gt;
            </div>
        </div>

        <div class="demo-section">
            <h2>📚 Документация</h2>
            <p>Больше информации вы найдете в README.md и USAGE.md</p>
            <a href="/" class="button">Главная страница</a>
            <a href="/demo-custom" class="button">Кастомная тема</a>
            <a href="/demo-prompt" class="button">Программист-помощник</a>
        </div>

        <div class="footer">
            <p>AI Bot - Простой и удобный чат с AI</p>
            <p style="font-size: 14px; margin-top: 10px;">Попробуйте чат в правом нижнем углу! 👇</p>
        </div>
    </div>

    <!-- AI Bot - подключен одним тегом! -->
    <script src="` + baseURL + `/chat.js"></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, tmpl)
}

func serveCustomDemoPage(w http.ResponseWriter, r *http.Request) {
	// Получаем базовый URL для скрипта
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	tmpl := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Bot - Кастомная тема</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
        }
        .header {
            text-align: center;
            color: white;
            margin-bottom: 50px;
        }
        .header h1 {
            font-size: 48px;
            margin-bottom: 10px;
            text-shadow: 0 2px 10px rgba(0,0,0,0.2);
        }
        .header p {
            font-size: 20px;
            opacity: 0.9;
        }
        .demo-section {
            background: white;
            border-radius: 20px;
            padding: 40px;
            margin-bottom: 30px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
        }
        .demo-section h2 {
            color: #e74c3c;
            margin-bottom: 20px;
            font-size: 28px;
        }
        .demo-section p {
            color: #666;
            line-height: 1.6;
            margin-bottom: 15px;
            font-size: 16px;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
            color: white;
            text-decoration: none;
            border-radius: 25px;
            font-weight: 600;
            margin: 10px 10px 10px 0;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        .button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(231, 76, 60, 0.4);
        }
        .footer {
            text-align: center;
            color: white;
            margin-top: 50px;
            opacity: 0.8;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔥 AI Bot</h1>
            <p>Демо с кастомной красной темой</p>
        </div>

        <div class="demo-section">
            <h2>🎨 Кастомная тема</h2>
            <p>Этот чат использует красную цветовую схему вместо стандартной сине-фиолетовой.</p>
            <p>Цвета настраиваются через data-атрибуты в теге script.</p>
        </div>

        <div class="demo-section">
            <h2>🚀 Возможности</h2>
            <p>• Перетаскивание кнопки и окна чата мышкой</p>
            <p>• Сохранение позиции в localStorage</p>
            <p>• Кастомизация цветов через data-атрибуты</p>
            <p>• Дополнительные CSS стили</p>
            <p>• Адаптивный дизайн для мобильных</p>
        </div>

        <div class="demo-section">
            <h2>📚 Навигация</h2>
            <a href="/" class="button">Главная страница</a>
            <a href="/demo" class="button">Обычная демо</a>
        </div>

        <div class="footer">
            <p>AI Bot - Попробуйте перетащить чат! 👇</p>
        </div>
    </div>

    <!-- AI Bot с кастомными красными цветами -->
    <script src="` + baseURL + `/chat.js"
            data-primary-color="#e74c3c"
            data-secondary-color="#c0392b"
            data-accent-color="#f39c12"></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, tmpl)
}

func servePromptDemoPage(w http.ResponseWriter, r *http.Request) {
	// Получаем базовый URL для скрипта
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	tmpl := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Bot - Программист-помощник</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Courier New', monospace;
            background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
            min-height: 100vh;
            padding: 40px 20px;
            color: #ecf0f1;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
        }
        .header {
            text-align: center;
            color: #ecf0f1;
            margin-bottom: 50px;
        }
        .header h1 {
            font-size: 48px;
            margin-bottom: 10px;
            text-shadow: 0 2px 10px rgba(0,0,0,0.3);
        }
        .header p {
            font-size: 20px;
            opacity: 0.9;
        }
        .demo-section {
            background: rgba(52, 73, 94, 0.8);
            border-radius: 10px;
            padding: 30px;
            margin-bottom: 30px;
            border: 1px solid #34495e;
        }
        .demo-section h2 {
            color: #3498db;
            margin-bottom: 20px;
            font-size: 24px;
        }
        .demo-section p {
            color: #bdc3c7;
            line-height: 1.6;
            margin-bottom: 15px;
            font-size: 16px;
        }
        .code-block {
            background: #1e1e1e;
            border: 1px solid #444;
            border-radius: 5px;
            padding: 20px;
            margin: 20px 0;
            overflow-x: auto;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            color: #f8f8f2;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
            color: white;
            text-decoration: none;
            border-radius: 5px;
            font-weight: 600;
            margin: 10px 10px 10px 0;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        .button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(52, 152, 219, 0.4);
        }
        .footer {
            text-align: center;
            color: #bdc3c7;
            margin-top: 50px;
            opacity: 0.8;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💻 AI Bot</h1>
            <p>Демо: Программист-помощник</p>
        </div>

        <div class="demo-section">
            <h2>🤖 Кастомный системный промпт</h2>
            <p>Этот чат настроен как дружелюбный помощник программиста.</p>
            <p>Он будет объяснять код простыми словами, давать советы и использовать эмодзи.</p>
        </div>

        <div class="demo-section">
            <h2>💡 Попробуйте спросить:</h2>
            <p>• "Объясни что такое замыкания в JavaScript"</p>
            <p>• "Как работает async/await?"</p>
            <p>• "Покажи пример REST API на Go"</p>
            <p>• "В чем разница между let и const?"</p>
        </div>

        <div class="demo-section">
            <h2>⚙️ Настройка промпта</h2>
            <div class="code-block">
data-system-prompt="Ты дружелюбный помощник программиста. 
Объясняй код простыми словами, давай практические советы, 
используй эмодзи для наглядности. Будь терпеливым и 
поддерживающим, особенно с новичками."
            </div>
        </div>

        <div class="demo-section">
            <h2>📚 Навигация</h2>
            <a href="/" class="button">Главная</a>
            <a href="/demo" class="button">Обычная демо</a>
            <a href="/demo-custom" class="button">Кастомная тема</a>
        </div>

        <div class="footer">
            <p>AI Bot - Ваш персональный помощник в программировании! 🚀</p>
        </div>
    </div>

    <!-- AI Bot с кастомным промптом для программистов -->
    <script src="` + baseURL + `/chat.js"
            data-primary-color="#3498db"
            data-secondary-color="#2980b9"
            data-accent-color="#e74c3c"
            data-system-prompt="Ты дружелюбный помощник программиста. Объясняй код простыми словами, давай практические советы, используй эмодзи для наглядности. Будь терпеливым и поддерживающим, особенно с новичками."></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, tmpl)
}

func handleChat(w http.ResponseWriter, r *http.Request, client *ai.Client, aiConfig *ai.Config) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message      string           `json:"message"`
		History      []ai.ChatMessage `json:"history"`
		SystemPrompt string           `json:"systemPrompt,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получаем системный промпт из запроса или конфигурации
	systemPrompt := req.SystemPrompt
	if systemPrompt == "" {
		cfg, _ := config.Load()
		systemPrompt = cfg.SystemPrompt
		if systemPrompt == "" {
			systemPrompt = "Ты полезный AI ассистент. Отвечай кратко и по делу на русском языке."
		}
	}

	// Создаем сообщения для AI
	messages := []ai.ChatMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// Добавляем историю
	messages = append(messages, req.History...)

	// Добавляем текущее сообщение
	messages = append(messages, ai.ChatMessage{
		Role:    "user",
		Content: req.Message,
	})

	// Отправляем запрос к AI
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(aiConfig.RequestTimeout)*time.Second)
	defer cancel()

	response, err := client.Chat(ctx, messages)
	if err != nil {
		http.Error(w, fmt.Sprintf("AI error: %v", err), http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	})
}

func handleStatus(w http.ResponseWriter, r *http.Request, client *ai.Client) {
	w.Header().Set("Content-Type", "application/json")

	status := map[string]interface{}{
		"configured": client.IsConfigured(),
		"provider":   client.GetProvider(),
	}

	// Проверяем доступность
	if client.IsConfigured() {
		testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		testMessages := []ai.ChatMessage{
			{Role: "user", Content: "test"},
		}

		_, err := client.Chat(testCtx, testMessages)
		status["available"] = err == nil
		if err != nil {
			status["error"] = err.Error()
		}
	} else {
		status["available"] = false
		status["error"] = "AI не настроен"
	}

	json.NewEncoder(w).Encode(status)
}

func runConfig() {
	fmt.Println("🔧 Конфигурация AI Bot")
	fmt.Println("======================")
	fmt.Println()

	// Загружаем текущую конфигурацию
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		return
	}

	// Проверяем наличие OpenRouter API ключа
	if cfg.OpenRouterKey == "" {
		fmt.Println("⚠️  OpenRouter API ключ не найден в .env файле")
		fmt.Print("Введите OpenRouter API ключ: ")
		reader := bufio.NewReader(os.Stdin)
		key, _ := reader.ReadString('\n')
		cfg.OpenRouterKey = strings.TrimSpace(key)

		if cfg.OpenRouterKey == "" {
			fmt.Println("❌ API ключ не может быть пустым")
			return
		}
	}

	fmt.Println("📡 Загрузка моделей из OpenRouter...")
	fmt.Println()

	// Создаем временный AI клиент для получения моделей
	aiConfig := &ai.Config{
		OpenRouterAPIKey: cfg.OpenRouterKey,
		OpenRouterURL:    "https://openrouter.ai/api/v1",
		RequestTimeout:   30,
	}
	client := ai.NewClient(aiConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	models, err := client.GetModels(ctx)
	if err != nil {
		fmt.Printf("❌ Ошибка загрузки моделей: %v\n", err)
		return
	}

	if len(models) == 0 {
		fmt.Println("❌ Модели не найдены")
		return
	}

	// Запускаем интерактивный выбор модели
	selectedModel := selectModelInteractive(models)
	if selectedModel == "" {
		fmt.Println("❌ Модель не выбрана")
		return
	}

	// Сохраняем выбранную модель
	cfg.OpenRouterModel = selectedModel

	// Настройка системного промпта
	fmt.Println()
	fmt.Println("🎯 Настройка системного промпта")
	fmt.Println("================================")
	fmt.Printf("Текущий промпт: %s\n", cfg.SystemPrompt)
	fmt.Println()
	fmt.Print("Хотите изменить системный промпт? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" || answer == "да" {
		fmt.Println()
		fmt.Println("Примеры системных промптов:")
		fmt.Println("1. Дружелюбный помощник:")
		fmt.Println("   Ты дружелюбный AI помощник. Общайся тепло и неформально, используй эмодзи.")
		fmt.Println()
		fmt.Println("2. Профессиональный консультант:")
		fmt.Println("   Ты профессиональный консультант. Давай четкие, структурированные ответы с примерами.")
		fmt.Println()
		fmt.Println("3. Программист-наставник:")
		fmt.Println("   Ты опытный программист. Помогай с кодом, объясняй концепции, предлагай лучшие практики.")
		fmt.Println()
		fmt.Println("4. Креативный помощник:")
		fmt.Println("   Ты креативный помощник. Генерируй идеи, помогай с творческими задачами, вдохновляй.")
		fmt.Println()
		fmt.Print("Введите новый системный промпт (или нажмите Enter для пропуска): ")
		
		newPrompt, _ := reader.ReadString('\n')
		newPrompt = strings.TrimSpace(newPrompt)
		
		if newPrompt != "" {
			cfg.SystemPrompt = newPrompt
			fmt.Printf("✅ Системный промпт обновлен\n")
		}
	}

	// Сохраняем конфигурацию
	if err := config.Save(cfg); err != nil {
		fmt.Printf("❌ Ошибка сохранения конфигурации: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Printf("✅ Конфигурация успешно сохранена в .env файл\n")
	fmt.Printf("   Модель: %s\n", selectedModel)
	fmt.Printf("   Промпт: %s\n", cfg.SystemPrompt)
	fmt.Println()
	fmt.Println("Теперь вы можете запустить бота:")
	fmt.Println("  ./ai-bot.exe")
}

// TUI модель для выбора AI модели
type modelSelectorModel struct {
	categories      []modelCategory
	currentCategory int
	currentModel    int
	selectedModel   string
	quitting        bool
	width           int
	height          int
}

type modelCategory struct {
	name   string
	models []ai.ModelInfo
	icon   string
}

func (m modelSelectorModel) Init() tea.Cmd {
	return nil
}

func (m modelSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "left", "h":
			m.currentCategory--
			if m.currentCategory < 0 {
				m.currentCategory = len(m.categories) - 1
			}
			m.currentModel = 0

		case "right", "l":
			m.currentCategory++
			if m.currentCategory >= len(m.categories) {
				m.currentCategory = 0
			}
			m.currentModel = 0

		case "up", "k":
			if len(m.categories[m.currentCategory].models) > 0 {
				m.currentModel--
				if m.currentModel < 0 {
					m.currentModel = len(m.categories[m.currentCategory].models) - 1
				}
			}

		case "down", "j":
			if len(m.categories[m.currentCategory].models) > 0 {
				m.currentModel++
				if m.currentModel >= len(m.categories[m.currentCategory].models) {
					m.currentModel = 0
				}
			}

		case "enter", " ":
			if len(m.categories[m.currentCategory].models) > 0 {
				m.selectedModel = m.categories[m.currentCategory].models[m.currentModel].ID
				m.quitting = true
				return m, tea.Quit
			}

		case "1":
			m.currentCategory = 0
			m.currentModel = 0
		case "2":
			if len(m.categories) > 1 {
				m.currentCategory = 1
				m.currentModel = 0
			}
		case "3":
			if len(m.categories) > 2 {
				m.currentCategory = 2
				m.currentModel = 0
			}
		}
	}

	return m, nil
}

func (m modelSelectorModel) View() string {
	if m.quitting {
		return ""
	}

	// Стили
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	tabStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(0, 1)

	activeTabStyle := tabStyle.Copy().
		Bold(true).
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#874BFD"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#EE6FF8")).
		Background(lipgloss.Color("#2A2A2A")).
		Padding(0, 1)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		MarginLeft(2)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		MarginTop(1)

	var s strings.Builder

	// Заголовок
	s.WriteString(titleStyle.Render("🤖 AI Bot - Выбор модели"))
	s.WriteString("\n\n")

	// Вкладки категорий
	var tabs []string
	for i, cat := range m.categories {
		style := tabStyle
		if i == m.currentCategory {
			style = activeTabStyle
		}
		tabs = append(tabs, style.Render(fmt.Sprintf("%s %s (%d)", cat.icon, cat.name, len(cat.models))))
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	s.WriteString("\n\n")

	// Список моделей
	currentModels := m.categories[m.currentCategory].models
	if len(currentModels) == 0 {
		s.WriteString("Модели не найдены в этой категории\n")
	} else {
		// Показываем модели с прокруткой
		maxVisible := 12
		startIdx := 0
		endIdx := len(currentModels)

		if len(currentModels) > maxVisible {
			if m.currentModel >= maxVisible/2 {
				startIdx = m.currentModel - maxVisible/2
				if startIdx+maxVisible > len(currentModels) {
					startIdx = len(currentModels) - maxVisible
				}
			}
			endIdx = startIdx + maxVisible
			if endIdx > len(currentModels) {
				endIdx = len(currentModels)
			}
		}

		for i := startIdx; i < endIdx; i++ {
			model := currentModels[i]
			cursor := " "
			if i == m.currentModel {
				cursor = ">"
				s.WriteString(selectedStyle.Render(fmt.Sprintf("%s %s", cursor, model.ID)))
			} else {
				s.WriteString(fmt.Sprintf("%s %s", cursor, model.ID))
			}
			s.WriteString("\n")

			// Показываем описание для выбранной модели
			if i == m.currentModel {
				if model.Name != "" {
					s.WriteString(descStyle.Render(fmt.Sprintf("📝 %s", model.Name)))
					s.WriteString("\n")
				}
				if model.Description != "" {
					desc := model.Description
					if len(desc) > 80 {
						desc = desc[:77] + "..."
					}
					s.WriteString(descStyle.Render(fmt.Sprintf("💬 %s", desc)))
					s.WriteString("\n")
				}
			}
		}

		// Индикатор прокрутки
		if len(currentModels) > maxVisible {
			s.WriteString(fmt.Sprintf("\nПоказано %d-%d из %d моделей", startIdx+1, endIdx, len(currentModels)))
			s.WriteString("\n")
		}
	}

	// Помощь
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("Управление: ↑↓/jk - навигация, ←→/hl - категории, Enter/Space - выбор, q/Esc - выход"))

	return s.String()
}

// Интерактивный выбор модели с TUI интерфейсом
func selectModelInteractive(models []ai.ModelInfo) string {
	// Категоризируем модели
	freeModels := []ai.ModelInfo{}
	popularModels := []ai.ModelInfo{}
	allModels := []ai.ModelInfo{}

	popularIDs := map[string]bool{
		"anthropic/claude-3.5-sonnet":       true,
		"anthropic/claude-3-haiku":          true,
		"openai/gpt-4o":                     true,
		"openai/gpt-4o-mini":                true,
		"meta-llama/llama-3.1-70b-instruct": true,
		"meta-llama/llama-3.1-8b-instruct":  true,
		"google/gemini-pro-1.5":             true,
		"mistralai/mistral-7b-instruct":     true,
	}

	// Определяем бесплатные модели (по ключевым словам)
	freeKeywords := []string{"free", "llama", "mistral", "qwen", "phi", "gemma", "deepseek"}

	for _, model := range models {
		allModels = append(allModels, model)

		// Проверяем на бесплатность
		isFree := false
		modelLower := strings.ToLower(model.ID)
		for _, keyword := range freeKeywords {
			if strings.Contains(modelLower, keyword) {
				isFree = true
				break
			}
		}
		if isFree {
			freeModels = append(freeModels, model)
		}

		// Проверяем популярность
		if popularIDs[model.ID] {
			popularModels = append(popularModels, model)
		}
	}

	// Сортируем модели
	sort.Slice(freeModels, func(i, j int) bool { return freeModels[i].ID < freeModels[j].ID })
	sort.Slice(popularModels, func(i, j int) bool { return popularModels[i].ID < popularModels[j].ID })
	sort.Slice(allModels, func(i, j int) bool { return allModels[i].ID < allModels[j].ID })

	// Создаем модель TUI
	m := modelSelectorModel{
		categories: []modelCategory{
			{"Бесплатные", freeModels, "🆓"},
			{"Популярные", popularModels, "⭐"},
			{"Все модели", allModels, "📋"},
		},
		currentCategory: 0,
		currentModel:    0,
	}

	// Запускаем TUI
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Ошибка TUI: %v\n", err)
		return ""
	}

	// Возвращаем выбранную модель
	if result, ok := finalModel.(modelSelectorModel); ok {
		return result.selectedModel
	}

	return ""
}