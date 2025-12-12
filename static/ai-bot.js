// AI Bot JavaScript клиент
class AIBot {
    constructor() {
        this.apiUrl = '/api/chat';
        this.statusUrl = '/api/status';
        this.history = [];
        this.init();
    }

    init() {
        this.chatContainer = document.getElementById('chatContainer');
        this.messageInput = document.getElementById('messageInput');
        this.sendButton = document.getElementById('sendButton');
        this.statusElement = document.getElementById('status');

        // Обработчики событий
        this.sendButton.addEventListener('click', () => this.sendMessage());
        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });

        // Проверка статуса при загрузке
        this.checkStatus();
    }

    async checkStatus() {
        try {
            const response = await fetch(this.statusUrl);
            const status = await response.json();
            
            if (status.available) {
                this.statusElement.textContent = `✓ Подключено (${status.provider})`;
                this.statusElement.style.color = '#4caf50';
            } else {
                this.statusElement.textContent = `✗ Ошибка: ${status.error || 'Неизвестная ошибка'}`;
                this.statusElement.style.color = '#f44336';
            }
        } catch (error) {
            this.statusElement.textContent = '✗ Ошибка подключения';
            this.statusElement.style.color = '#f44336';
        }
    }

    async sendMessage() {
        const message = this.messageInput.value.trim();
        if (!message) return;

        // Добавляем сообщение пользователя в чат
        this.addMessage('user', message);
        this.messageInput.value = '';
        
        // Блокируем кнопку и поле ввода
        this.setLoading(true);

        try {
            const response = await fetch(this.apiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    message: message,
                    history: this.history,
                }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            
            // Добавляем сообщение пользователя в историю
            this.history.push({
                role: 'user',
                content: message,
            });

            // Добавляем ответ бота в историю и чат
            this.history.push({
                role: 'assistant',
                content: data.response,
            });

            this.addMessage('bot', data.response);
        } catch (error) {
            console.error('Error:', error);
            this.addMessage('bot', `Ошибка: ${error.message}. Проверьте конфигурацию API ключей.`);
        } finally {
            this.setLoading(false);
        }
    }

    addMessage(role, content) {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${role}`;
        
        const contentDiv = document.createElement('div');
        contentDiv.className = 'message-content';
        contentDiv.textContent = content;
        
        messageDiv.appendChild(contentDiv);
        this.chatContainer.appendChild(messageDiv);
        
        // Прокрутка вниз
        this.chatContainer.scrollTop = this.chatContainer.scrollHeight;
    }

    setLoading(loading) {
        this.sendButton.disabled = loading;
        this.messageInput.disabled = loading;
        
        if (loading) {
            this.sendButton.textContent = 'Отправка...';
        } else {
            this.sendButton.textContent = 'Отправить';
        }
    }
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    new AIBot();
});

