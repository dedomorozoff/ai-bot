package main

import (
	"fmt"
	"net/http"
)

func serveEmbeddedChatNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "no-cache")

	// Получаем базовый URL для API
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	// Создаем JavaScript код
	js := `(function() {
	// Получаем кастомные CSS из data-атрибутов скрипта
	var scriptTag = document.currentScript || document.querySelector('script[src*="chat.js"]');
	var customCSS = '';
	var customColors = {};
	
	var systemPrompt = '';
	
	if (scriptTag) {
		// Читаем data-атрибуты для кастомизации
		customColors.primary = scriptTag.getAttribute('data-primary-color') || '#667eea';
		customColors.secondary = scriptTag.getAttribute('data-secondary-color') || '#764ba2';
		customColors.accent = scriptTag.getAttribute('data-accent-color') || '#ff4757';
		customCSS = scriptTag.getAttribute('data-custom-css') || '';
		systemPrompt = scriptTag.getAttribute('data-system-prompt') || '';
	}

	// Базовые CSS стили для AI виджета
	var css = '.ai-chat-widget{position:fixed;bottom:20px;right:20px;z-index:10000;font-family:Inter,-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif;user-select:none}.ai-chat-toggle{width:60px;height:60px;background:linear-gradient(135deg,' + customColors.primary + ' 0%,' + customColors.secondary + ' 100%);border-radius:50%;display:flex;align-items:center;justify-content:center;cursor:pointer;box-shadow:0 4px 12px rgba(0,0,0,.3);transition:all .3s ease;position:relative;border:none}.ai-chat-toggle:hover{transform:scale(1.05);box-shadow:0 6px 20px rgba(0,0,0,.4)}.ai-chat-toggle.dragging{transform:scale(1.1);box-shadow:0 8px 25px rgba(0,0,0,.5);transition:none}.ai-chat-toggle-icon{color:white;font-size:24px;pointer-events:none}.ai-chat-badge{position:absolute;top:-5px;right:-5px;background:' + customColors.accent + ';color:white;font-size:10px;padding:2px 6px;border-radius:10px;font-weight:bold;pointer-events:none}.ai-chat-window{position:absolute;bottom:80px;right:0;width:350px;height:500px;background:white;border-radius:12px;box-shadow:0 8px 30px rgba(0,0,0,.3);display:none;flex-direction:column;overflow:hidden;border:1px solid #e0e0e0}.ai-chat-window.open{display:flex;animation:slideUp .3s ease}@keyframes slideUp{from{opacity:0;transform:translateY(20px)}to{opacity:1;transform:translateY(0)}}.ai-chat-header{background:linear-gradient(135deg,' + customColors.primary + ' 0%,' + customColors.secondary + ' 100%);color:white;padding:15px;display:flex;justify-content:space-between;align-items:center;cursor:move}.ai-chat-header.dragging{cursor:grabbing}.ai-chat-title{display:flex;align-items:center;gap:8px;font-weight:600;pointer-events:none}.ai-chat-close{background:none;border:none;color:white;cursor:pointer;padding:4px;border-radius:4px;transition:background .2s;font-size:16px}.ai-chat-close:hover{background:rgba(255,255,255,.2)}.ai-chat-messages{flex:1;padding:15px;overflow-y:auto;background:#f8f9fa}.ai-message,.user-message{display:flex;margin-bottom:15px;align-items:flex-start;gap:10px}.user-message{flex-direction:row-reverse}.ai-avatar,.user-avatar{width:32px;height:32px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:16px;flex-shrink:0}.ai-avatar{background:#e9ecef}.user-avatar{background:linear-gradient(135deg,' + customColors.primary + ' 0%,' + customColors.secondary + ' 100%);color:white}.ai-message-content,.user-message-content{max-width:80%}.ai-message-text,.user-message-text{background:white;color:#333;padding:10px 12px;border-radius:12px;border:1px solid #e0e0e0;line-height:1.4}.user-message-text{background:linear-gradient(135deg,' + customColors.primary + ' 0%,' + customColors.secondary + ' 100%);color:white;border:none}.ai-message-time,.user-message-time{font-size:11px;color:#6c757d;margin-top:4px;padding:0 4px}.ai-quick-buttons{padding:10px 15px;display:flex;gap:6px;flex-wrap:wrap;background:white}.ai-quick-btn{background:#f8f9fa;color:#495057;border:1px solid #dee2e6;padding:6px 10px;border-radius:16px;font-size:12px;cursor:pointer;transition:all .2s}.ai-quick-btn:hover{background:#e9ecef}.ai-input-row{display:flex;padding:15px;gap:10px;background:white;border-top:1px solid #e0e0e0}.ai-input-row input{flex:1;border:1px solid #dee2e6;background:white;border-radius:20px;padding:10px 15px;font-size:14px;outline:none;transition:border-color .2s}.ai-input-row input:focus{border-color:' + customColors.primary + '}.ai-input-row input::placeholder{color:#6c757d}.ai-input-row button{width:40px;height:40px;background:linear-gradient(135deg,' + customColors.primary + ' 0%,' + customColors.secondary + ' 100%);border:none;border-radius:50%;color:white;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:opacity .2s;font-size:16px}.ai-input-row button:hover{opacity:.9}.ai-input-row button:disabled{background:#6c757d;cursor:not-allowed}.ai-typing{display:flex;align-items:center;gap:4px;padding:8px 12px;background:white;border:1px solid #e0e0e0;border-radius:12px}.ai-typing-dot{width:6px;height:6px;background:#6c757d;border-radius:50%;animation:typing 1.4s infinite ease-in-out}.ai-typing-dot:nth-child(1){animation-delay:-.32s}.ai-typing-dot:nth-child(2){animation-delay:-.16s}.ai-typing-dot:nth-child(3){animation-delay:0s}@keyframes typing{0%,80%,100%{transform:scale(.8);opacity:.5}40%{transform:scale(1);opacity:1}}@media (max-width:768px){.ai-chat-window{width:300px;height:450px}.ai-chat-widget{bottom:15px;right:15px}}';
	
	// Добавляем кастомные CSS если есть
	if (customCSS) {
		css += customCSS;
	}
	
	// Добавляем стили
	var style = document.createElement('style');
	style.textContent = css;
	document.head.appendChild(style);

	var widget = null;
	var chatWindow = null;
	var messages = null;
	var input = null;
	var sendBtn = null;
	var badge = null;
	var toggle = null;
	var header = null;
	var history = [];
	var isOpen = false;
	var isTyping = false;
	var isDragging = false;
	var dragTarget = null;
	var dragOffset = {x: 0, y: 0};
	var apiUrl = '` + baseURL + `/api/chat';
	
	function initChat() {
		if (widget) return;
		
		widget = document.createElement('div');
		widget.className = 'ai-chat-widget';
		widget.innerHTML = '<button class="ai-chat-toggle" onclick="toggleAIChat()"><span class="ai-chat-toggle-icon">🤖</span><span class="ai-chat-badge" id="aiChatBadge">AI</span></button><div class="ai-chat-window" id="aiChatWindow"><div class="ai-chat-header"><div class="ai-chat-title"><span>🤖</span><span>AI Помощник</span></div><button class="ai-chat-close" onclick="closeAIChat()">×</button></div><div class="ai-chat-messages" id="aiChatMessages"><div class="ai-message"><div class="ai-avatar">🤖</div><div class="ai-message-content"><div class="ai-message-text">Привет! Я помогу вам с любыми вопросами. Просто напишите что вас интересует.</div><div class="ai-message-time">сейчас</div></div></div></div><div class="ai-quick-buttons"><button class="ai-quick-btn" onclick="sendQuickMessage(\'Как дела?\')">👋 Привет</button><button class="ai-quick-btn" onclick="sendQuickMessage(\'Помоги с кодом\')">💻 Код</button><button class="ai-quick-btn" onclick="sendQuickMessage(\'Объясни концепцию\')">📚 Обучение</button></div><div class="ai-input-row"><input type="text" id="aiChatInput" placeholder="Напишите сообщение..." maxlength="500"><button id="aiSendButton" onclick="sendAIMessage()">➤</button></div></div>';
		
		document.body.appendChild(widget);
		
		toggle = widget.querySelector('.ai-chat-toggle');
		chatWindow = document.getElementById('aiChatWindow');
		header = chatWindow.querySelector('.ai-chat-header');
		messages = document.getElementById('aiChatMessages');
		input = document.getElementById('aiChatInput');
		sendBtn = document.getElementById('aiSendButton');
		badge = document.getElementById('aiChatBadge');
		
		input.addEventListener('keypress', function(e) {
			if (e.key === 'Enter' && !e.shiftKey) {
				e.preventDefault();
				sendMessage();
			}
		});

		document.addEventListener('click', function(e) {
			if (isOpen && !widget.contains(e.target)) {
				closeChat();
			}
		});

		// Добавляем обработчики перетаскивания
		setupDragHandlers();
		
		checkAIStatus();
	}

	window.toggleAIChat = function() {
		if (isOpen) {
			closeChat();
		} else {
			openChat();
		}
	};

	window.closeAIChat = function() {
		closeChat();
	};

	window.sendAIMessage = function() {
		sendMessage();
	};

	window.sendQuickMessage = function(message) {
		input.value = message;
		sendMessage();
	};

	function openChat() {
		isOpen = true;
		chatWindow.classList.add('open');
		input.focus();
		badge.style.display = 'none';
	}

	function closeChat() {
		isOpen = false;
		chatWindow.classList.remove('open');
		badge.style.display = 'block';
	}

	async function checkAIStatus() {
		try {
			var response = await fetch('` + baseURL + `/api/status');
			var status = await response.json();
			
			if (status.configured && status.available) {
				badge.textContent = 'AI';
				badge.style.background = '#2ed573';
			} else {
				badge.textContent = '!';
				badge.style.background = '#ff4757';
			}
		} catch (error) {
			badge.textContent = '?';
			badge.style.background = '#ffa502';
		}
	}

	async function sendMessage() {
		var message = input.value.trim();
		if (!message || isTyping) return;

		addMessage(message, 'user');
		input.value = '';
		showTyping();

		history.push({role: 'user', content: message});

		try {
			var requestBody = {message: message, history: history};
			if (systemPrompt) {
				requestBody.systemPrompt = systemPrompt;
			}
			
			var response = await fetch(apiUrl, {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify(requestBody)
			});

			if (!response.ok) throw new Error('HTTP ' + response.status);

			var data = await response.json();
			
			hideTyping();
			addMessage(data.response, 'ai');
			history.push({role: 'assistant', content: data.response});

		} catch (error) {
			hideTyping();
			addMessage('Извините, произошла ошибка. Попробуйте еще раз.', 'ai');
		}
	}

	function addMessage(content, sender) {
		var messageDiv = document.createElement('div');
		messageDiv.className = sender + '-message';
		
		var now = new Date().toLocaleTimeString('ru-RU', { 
			hour: '2-digit', 
			minute: '2-digit' 
		});
		
		var avatar = sender === 'user' ? '👤' : '🤖';
		
		messageDiv.innerHTML = '<div class="' + sender + '-avatar">' + avatar + '</div><div class="' + sender + '-message-content"><div class="' + sender + '-message-text">' + formatMessage(content) + '</div><div class="' + sender + '-message-time">' + now + '</div></div>';
		
		messages.appendChild(messageDiv);
		scrollToBottom();
	}

	function formatMessage(content) {
		return content
			.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
			.replace(/\*(.*?)\*/g, '<em>$1</em>')
			.replace(/` + "`" + `(.*?)` + "`" + `/g, '<code style="background:#f1f3f4;padding:2px 4px;border-radius:3px;">$1</code>')
			.replace(/\n/g, '<br>');
	}

	function showTyping() {
		isTyping = true;
		sendBtn.disabled = true;
		
		var typingDiv = document.createElement('div');
		typingDiv.className = 'ai-message';
		typingDiv.id = 'typingIndicator';
		typingDiv.innerHTML = '<div class="ai-avatar">🤖</div><div class="ai-message-content"><div class="ai-typing"><div class="ai-typing-dot"></div><div class="ai-typing-dot"></div><div class="ai-typing-dot"></div></div></div>';
		
		messages.appendChild(typingDiv);
		scrollToBottom();
	}

	function hideTyping() {
		isTyping = false;
		sendBtn.disabled = false;
		
		var typingIndicator = document.getElementById('typingIndicator');
		if (typingIndicator) {
			typingIndicator.remove();
		}
	}

	function scrollToBottom() {
		messages.scrollTop = messages.scrollHeight;
	}

	function setupDragHandlers() {
		// Перетаскивание кнопки
		toggle.addEventListener('mousedown', function(e) {
			if (e.button !== 0) return; // Только левая кнопка мыши
			startDrag(e, 'toggle');
		});

		// Перетаскивание окна чата за заголовок
		header.addEventListener('mousedown', function(e) {
			if (e.button !== 0) return; // Только левая кнопка мыши
			if (e.target.classList.contains('ai-chat-close')) return; // Не перетаскиваем при клике на кнопку закрытия
			startDrag(e, 'window');
		});

		// Глобальные обработчики
		document.addEventListener('mousemove', handleDrag);
		document.addEventListener('mouseup', stopDrag);
		
		// Предотвращаем выделение текста при перетаскивании
		document.addEventListener('selectstart', function(e) {
			if (isDragging) e.preventDefault();
		});
	}

	function startDrag(e, target) {
		isDragging = true;
		dragTarget = target;
		
		var rect = widget.getBoundingClientRect();
		dragOffset.x = e.clientX - rect.left;
		dragOffset.y = e.clientY - rect.top;
		
		// Добавляем класс для визуального эффекта
		if (target === 'toggle') {
			toggle.classList.add('dragging');
		} else {
			header.classList.add('dragging');
		}
		
		e.preventDefault();
	}

	function handleDrag(e) {
		if (!isDragging) return;
		
		var newX = e.clientX - dragOffset.x;
		var newY = e.clientY - dragOffset.y;
		
		// Ограничиваем перемещение границами экрана
		var maxX = window.innerWidth - 60; // Ширина кнопки
		var maxY = window.innerHeight - 60; // Высота кнопки
		
		newX = Math.max(0, Math.min(newX, maxX));
		newY = Math.max(0, Math.min(newY, maxY));
		
		// Применяем новую позицию
		widget.style.left = newX + 'px';
		widget.style.top = newY + 'px';
		widget.style.right = 'auto';
		widget.style.bottom = 'auto';
		
		e.preventDefault();
	}

	function stopDrag(e) {
		if (!isDragging) return;
		
		isDragging = false;
		
		// Убираем классы перетаскивания
		toggle.classList.remove('dragging');
		header.classList.remove('dragging');
		
		// Сохраняем позицию в localStorage
		var rect = widget.getBoundingClientRect();
		localStorage.setItem('aiChatPosition', JSON.stringify({
			x: rect.left,
			y: rect.top
		}));
		
		dragTarget = null;
	}

	function loadSavedPosition() {
		try {
			var saved = localStorage.getItem('aiChatPosition');
			if (saved) {
				var pos = JSON.parse(saved);
				
				// Проверяем что позиция все еще в пределах экрана
				var maxX = window.innerWidth - 60;
				var maxY = window.innerHeight - 60;
				
				if (pos.x >= 0 && pos.x <= maxX && pos.y >= 0 && pos.y <= maxY) {
					widget.style.left = pos.x + 'px';
					widget.style.top = pos.y + 'px';
					widget.style.right = 'auto';
					widget.style.bottom = 'auto';
				}
			}
		} catch (e) {
			// Игнорируем ошибки загрузки позиции
		}
	}
	
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', function() {
			initChat();
			// Загружаем сохраненную позицию после инициализации
			setTimeout(loadSavedPosition, 100);
		});
	} else {
		initChat();
		setTimeout(loadSavedPosition, 100);
	}
})();`

	fmt.Fprint(w, js)
}