# 📖 Руководство по использованию AI Bot

Подробное руководство по настройке и использованию AI Bot виджета.

## 🚀 Быстрый старт

### Шаг 1: Получение API ключа

#### OpenRouter (рекомендуется)
1. Перейдите на [openrouter.ai](https://openrouter.ai)
2. Зарегистрируйтесь или войдите в аккаунт
3. Перейдите в раздел "Keys" 
4. Создайте новый API ключ
5. Скопируйте ключ (начинается с `sk-or-v1-...`)

#### OpenAI (альтернатива)
1. Перейдите на [platform.openai.com](https://platform.openai.com)
2. Войдите в аккаунт
3. Перейдите в "API Keys"
4. Создайте новый ключ
5. Скопируйте ключ (начинается с `sk-...`)

### Шаг 2: Настройка через TUI конфигуратор

```bash
./ai-bot.exe --config
```

Конфигуратор проведет вас через:

1. **Ввод API ключа** - вставьте скопированный ключ
2. **Загрузка моделей** - автоматически получит список доступных моделей
3. **Выбор модели** - выберите из категорий:
   - 🆓 **Бесплатные** - Llama, Mistral, Qwen (бесплатно через OpenRouter)
   - ⭐ **Популярные** - Claude, GPT-4, Gemini (проверенные модели)
   - 📋 **Все модели** - полный список доступных моделей
4. **Настройка промпта** - выберите поведение AI или создайте свой промпт

### Шаг 3: Запуск сервера

```bash
./ai-bot.exe
```

Сервер запустится на `http://localhost:8080`

### Шаг 4: Тестирование

Откройте в браузере:
- `http://localhost:8080` - основной чат
- `http://localhost:8080/demo` - демо с документацией

## 🎨 Примеры интеграции

### Базовая интеграция

```html
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Мой сайт</title>
</head>
<body>
    <h1>Добро пожаловать на мой сайт!</h1>
    <p>Здесь может быть любой контент...</p>
    
    <!-- AI Bot - просто добавьте эту строку -->
    <script src="http://localhost:8080/chat.js"></script>
</body>
</html>
```

### Интернет-магазин

```html
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Интернет-магазин</title>
    <style>
        body { font-family: Arial, sans-serif; }
        .product { border: 1px solid #ddd; padding: 20px; margin: 10px; }
    </style>
</head>
<body>
    <header>
        <h1>🛍️ Мой магазин</h1>
        <nav>
            <a href="#catalog">Каталог</a>
            <a href="#delivery">Доставка</a>
            <a href="#contacts">Контакты</a>
        </nav>
    </header>
    
    <main>
        <div class="product">
            <h3>Товар 1</h3>
            <p>Описание товара...</p>
            <button>Купить</button>
        </div>
    </main>
    
    <!-- AI помощник для интернет-магазина -->
    <script src="http://localhost:8080/chat.js"
            data-primary-color="#e67e22"
            data-secondary-color="#d35400"
            data-system-prompt="Ты помощник интернет-магазина. Помогай клиентам с выбором товаров, отвечай на вопросы о доставке, оплате и возврате. Будь дружелюбным и профессиональным."></script>
</body>
</html>
```

### Образовательный сайт

```html
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Онлайн курсы</title>
</head>
<body>
    <h1>📚 Изучай программирование</h1>
    
    <section>
        <h2>Курсы JavaScript</h2>
        <p>Изучите современный JavaScript с нуля...</p>
    </section>
    
    <!-- AI наставник для обучения -->
    <script src="http://localhost:8080/chat.js"
            data-primary-color="#3498db"
            data-secondary-color="#2980b9"
            data-accent-color="#e74c3c"
            data-system-prompt="Ты опытный преподаватель программирования. Объясняй сложные концепции простыми словами, приводи примеры кода, помогай с домашними заданиями. Будь терпеливым и поддерживающим с новичками."></script>
</body>
</html>
```

### Корпоративный сайт

```html
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>IT Компания</title>
</head>
<body>
    <h1>🏢 Наша IT компания</h1>
    
    <section>
        <h2>Услуги</h2>
        <ul>
            <li>Разработка веб-приложений</li>
            <li>Мобильные приложения</li>
            <li>Консультации по IT</li>
        </ul>
    </section>
    
    <!-- Профессиональный консультант -->
    <script src="http://localhost:8080/chat.js"
            data-primary-color="#2c3e50"
            data-secondary-color="#34495e"
            data-system-prompt="Ты профессиональный IT консультант. Отвечай на вопросы о наших услугах, технологиях, сроках и стоимости проектов. Будь компетентным и деловым, но дружелюбным."
            data-custom-css=".ai-chat-window{border-radius:0;box-shadow:0 0 20px rgba(0,0,0,0.3);}"></script>
</body>
</html>
```

## 🎯 Настройка системных промптов

### Примеры готовых промптов

#### Дружелюбный помощник
```
Ты дружелюбный AI помощник. Общайся тепло и неформально, используй эмодзи 😊. 
Будь позитивным и поддерживающим в любых ситуациях. Помогай решать проблемы 
с энтузиазмом и всегда предлагай несколько вариантов решения.
```

#### Техническая поддержка
```
Ты специалист технической поддержки. Помогай пользователям решать технические 
проблемы пошагово. Задавай уточняющие вопросы, предлагай конкретные решения. 
Будь терпеливым и объясняй сложные вещи простыми словами.
```

#### Продавец-консультант
```
Ты опытный продавец-консультант. Помогай клиентам выбрать подходящий товар, 
рассказывай о преимуществах, отвечай на вопросы о характеристиках и доставке. 
Будь убедительным, но не навязчивым.
```

#### Преподаватель
```
Ты опытный преподаватель. Объясняй сложные концепции простыми словами, 
приводи примеры из жизни, проверяй понимание материала. Будь терпеливым 
с новичками и поощряй их успехи.
```

#### Креативный помощник
```
Ты креативный помощник и генератор идей. Помогай с творческими задачами, 
предлагай нестандартные решения, вдохновляй на новые проекты. Используй 
метафоры и аналогии для объяснения идей.
```

### Настройка через .env файл

```env
# Глобальный системный промпт для всех чатов
SYSTEM_PROMPT=Ты помощник нашей компании. Отвечай на вопросы о наших услугах, помогай клиентам, будь профессиональным и дружелюбным.
```

### Настройка через JavaScript

```html
<!-- Промпт только для этой страницы -->
<script src="http://localhost:8080/chat.js"
        data-system-prompt="Ты помощник интернет-магазина электроники. Знаешь все о смартфонах, ноутбуках и гаджетах. Помогай выбрать технику под конкретные задачи и бюджет."></script>
```

## 🎨 Кастомизация внешнего вида

### Цветовые схемы

#### Корпоративная (синяя)
```html
<script src="http://localhost:8080/chat.js"
        data-primary-color="#3498db"
        data-secondary-color="#2980b9"
        data-accent-color="#e74c3c"></script>
```

#### Природная (зеленая)
```html
<script src="http://localhost:8080/chat.js"
        data-primary-color="#27ae60"
        data-secondary-color="#229954"
        data-accent-color="#f39c12"></script>
```

#### Элегантная (фиолетовая)
```html
<script src="http://localhost:8080/chat.js"
        data-primary-color="#9b59b6"
        data-secondary-color="#8e44ad"
        data-accent-color="#e67e22"></script>
```

#### Энергичная (оранжевая)
```html
<script src="http://localhost:8080/chat.js"
        data-primary-color="#e67e22"
        data-secondary-color="#d35400"
        data-accent-color="#3498db"></script>
```

### Кастомные CSS стили

#### Квадратный дизайн
```html
<script src="http://localhost:8080/chat.js"
        data-custom-css=".ai-chat-toggle{border-radius:10px;}.ai-chat-window{border-radius:10px;}"></script>
```

#### Неоновая подсветка
```html
<script src="http://localhost:8080/chat.js"
        data-primary-color="#00ff88"
        data-secondary-color="#00cc6a"
        data-custom-css=".ai-chat-toggle{box-shadow:0 0 20px #00ff88;}.ai-chat-window{border:2px solid #00ff88;}"></script>
```

#### Минималистичный
```html
<script src="http://localhost:8080/chat.js"
        data-primary-color="#333"
        data-secondary-color="#555"
        data-custom-css=".ai-chat-toggle{box-shadow:none;border:1px solid #ddd;}.ai-chat-window{box-shadow:0 2px 10px rgba(0,0,0,0.1);}"></script>
```

#### Большая кнопка
```html
<script src="http://localhost:8080/chat.js"
        data-custom-css=".ai-chat-toggle{width:80px;height:80px;}.ai-chat-toggle-icon{font-size:32px;}"></script>
```

### Позиционирование

#### Левый нижний угол
```html
<script src="http://localhost:8080/chat.js"
        data-custom-css=".ai-chat-widget{left:20px;right:auto;}.ai-chat-window{left:0;right:auto;}"></script>
```

#### Верхний правый угол
```html
<script src="http://localhost:8080/chat.js"
        data-custom-css=".ai-chat-widget{top:20px;bottom:auto;}.ai-chat-window{top:80px;bottom:auto;}"></script>
```

#### По центру снизу
```html
<script src="http://localhost:8080/chat.js"
        data-custom-css=".ai-chat-widget{left:50%;right:auto;transform:translateX(-50%);}.ai-chat-window{left:50%;right:auto;transform:translateX(-50%);}"></script>
```

## 🔧 Продвинутые настройки

### Программное управление

```javascript
// Открыть чат программно
window.toggleAIChat();

// Закрыть чат
window.closeAIChat();

// Отправить сообщение программно
window.sendQuickMessage("Расскажи о ваших услугах");

// Прямой вызов API
async function askAI(question, customPrompt = null) {
    const body = {
        message: question,
        history: []
    };
    
    if (customPrompt) {
        body.systemPrompt = customPrompt;
    }
    
    const response = await fetch('/api/chat', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(body)
    });
    
    const data = await response.json();
    return data.response;
}

// Использование
const answer = await askAI("Что такое JavaScript?", "Объясни как учителю информатики");
console.log(answer);
```

### Интеграция с формами

```html
<form id="contactForm">
    <input type="text" id="name" placeholder="Ваше имя" required>
    <input type="email" id="email" placeholder="Email" required>
    <textarea id="message" placeholder="Сообщение" required></textarea>
    <button type="submit">Отправить</button>
    <button type="button" onclick="askAIForHelp()">Помощь AI</button>
</form>

<script>
async function askAIForHelp() {
    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const message = document.getElementById('message').value;
    
    let question = "Помоги заполнить форму обратной связи. ";
    if (!name) question += "Как лучше представиться? ";
    if (!email) question += "Какой email указать? ";
    if (!message) question += "Что написать в сообщении для IT компании? ";
    
    window.sendQuickMessage(question);
}
</script>
```

### Условная загрузка

```html
<script>
// Загружаем чат только для определенных страниц
if (window.location.pathname.includes('/support/') || 
    window.location.pathname.includes('/contact/')) {
    
    const script = document.createElement('script');
    script.src = 'http://localhost:8080/chat.js';
    script.setAttribute('data-primary-color', '#e74c3c');
    script.setAttribute('data-system-prompt', 'Ты специалист технической поддержки...');
    document.head.appendChild(script);
}
</script>
```

## 📱 Мобильная оптимизация

### Адаптивные стили

```html
<script src="http://localhost:8080/chat.js"
        data-custom-css="
        @media (max-width: 768px) {
            .ai-chat-widget { bottom: 10px; right: 10px; }
            .ai-chat-toggle { width: 50px; height: 50px; }
            .ai-chat-toggle-icon { font-size: 20px; }
            .ai-chat-window { width: calc(100vw - 20px); height: 70vh; }
        }"></script>
```

### Мобильное позиционирование

```html
<script src="http://localhost:8080/chat.js"
        data-custom-css="
        @media (max-width: 480px) {
            .ai-chat-widget { 
                left: 50%; 
                right: auto; 
                transform: translateX(-50%); 
                bottom: 10px; 
            }
            .ai-chat-window { 
                left: 10px; 
                right: 10px; 
                width: auto; 
                transform: none; 
            }
        }"></script>
```

## 🔍 Отладка и мониторинг

### Проверка статуса

```javascript
// Проверить доступность AI
fetch('/api/status')
    .then(r => r.json())
    .then(status => {
        console.log('AI Status:', status);
        if (!status.available) {
            console.error('AI недоступен:', status.error);
        }
    });
```

### Логирование сообщений

```javascript
// Перехватываем отправку сообщений для аналитики
const originalFetch = window.fetch;
window.fetch = function(...args) {
    if (args[0].includes('/api/chat')) {
        console.log('AI Chat Request:', args[1].body);
    }
    return originalFetch.apply(this, args);
};
```

### Обработка ошибок

```html
<script>
window.addEventListener('error', function(e) {
    if (e.message.includes('ai-chat')) {
        console.error('AI Chat Error:', e);
        // Отправить ошибку в систему мониторинга
        // analytics.track('ai_chat_error', {error: e.message});
    }
});
</script>
```

## 🚀 Развертывание в продакшене

### Nginx конфигурация

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Кэширование статических файлов
    location /chat.js {
        proxy_pass http://localhost:8080;
        expires 1h;
        add_header Cache-Control "public, immutable";
    }
}
```

### Docker Compose

```yaml
version: '3.8'
services:
  ai-bot:
    build: .
    ports:
      - "8080:8080"
    environment:
      - OPENROUTER_API_KEY=${OPENROUTER_API_KEY}
      - OPENROUTER_MODEL=anthropic/claude-3.5-sonnet
      - HOST=0.0.0.0
      - PORT=8080
    restart: unless-stopped
    volumes:
      - ./logs:/app/logs
```

### Systemd сервис

```ini
[Unit]
Description=AI Bot Service
After=network.target

[Service]
Type=simple
User=aibot
Group=aibot
WorkingDirectory=/opt/ai-bot
ExecStart=/opt/ai-bot/ai-bot
Restart=always
RestartSec=5
Environment=OPENROUTER_API_KEY=your_key_here

[Install]
WantedBy=multi-user.target
```

## 📊 Аналитика и метрики

### Отслеживание использования

```javascript
// Добавить в ваш сайт для отслеживания
let chatInteractions = 0;

// Перехватываем открытие чата
const originalToggle = window.toggleAIChat;
window.toggleAIChat = function() {
    chatInteractions++;
    // Отправить в Google Analytics
    gtag('event', 'ai_chat_toggle', {
        'event_category': 'engagement',
        'value': chatInteractions
    });
    return originalToggle.apply(this, arguments);
};
```

### Мониторинг производительности

```javascript
// Измерение времени ответа AI
const originalSendMessage = window.sendAIMessage;
window.sendAIMessage = function() {
    const startTime = performance.now();
    
    return originalSendMessage.apply(this, arguments).then(result => {
        const responseTime = performance.now() - startTime;
        console.log(`AI Response Time: ${responseTime}ms`);
        
        // Отправить метрику
        // analytics.track('ai_response_time', {duration: responseTime});
        
        return result;
    });
};
```

## 🔒 Безопасность

### Ограничение доступа

```javascript
// Проверка домена перед загрузкой чата
const allowedDomains = ['yourdomain.com', 'www.yourdomain.com'];
if (!allowedDomains.includes(window.location.hostname)) {
    console.warn('AI Chat не разрешен на этом домене');
} else {
    // Загружаем чат
    const script = document.createElement('script');
    script.src = 'http://localhost:8080/chat.js';
    document.head.appendChild(script);
}
```

### Валидация сообщений

```javascript
// Фильтрация нежелательного контента
function validateMessage(message) {
    const forbiddenWords = ['spam', 'hack', 'virus'];
    return !forbiddenWords.some(word => 
        message.toLowerCase().includes(word)
    );
}

// Использование в форме
document.getElementById('chatInput').addEventListener('input', function(e) {
    if (!validateMessage(e.target.value)) {
        e.target.style.borderColor = 'red';
    } else {
        e.target.style.borderColor = '';
    }
});
```

---

Это руководство покрывает основные сценарии использования AI Bot. Для получения дополнительной помощи обращайтесь к [документации](README.md) или создавайте [issue](https://github.com/username/ai-bot/issues).