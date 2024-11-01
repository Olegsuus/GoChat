// static/js/chat.js

document.addEventListener('DOMContentLoaded', () => {
    const path = window.location.pathname;
    if (!path.startsWith('/chats/')) return;

    const chatID = path.split('/')[2];
    const token = localStorage.getItem('token');
    const messagesDiv = document.getElementById('messages');
    const messageForm = document.getElementById('messageForm');
    const messageInput = document.getElementById('messageInput');
    const chatError = document.getElementById('chatError');

    if (!chatID || !token) return;

    const ws = new WebSocket(`ws://localhost:8765/api/chats/ws?chat_id=${chatID}`); // Замените PORT на ваш порт

    ws.onopen = () => {
        console.log('WebSocket подключен');
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        appendMessage(data.message);
    };

    ws.onerror = (error) => {
        console.error('WebSocket ошибка', error);
    };

    ws.onclose = () => {
        console.log('WebSocket отключен');
    };

    const appendMessage = (message) => {
        const messageElement = document.createElement('div');
        messageElement.className = 'message';
        messageElement.innerHTML = `
            <strong>${message.sender_id === '${userID}' ? 'Вы' : message.sender_id}:</strong>
            <span>${message.content}</span>
        `;
        messagesDiv.appendChild(messageElement);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
    };

    // Получение ID текущего пользователя из токена
    const userID = parseJwt(token)._id;

    function parseJwt(token) {
        try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
                return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
            }).join(''));
            return JSON.parse(jsonPayload);
        } catch (e) {
            return null;
        }
    }

    messageForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const content = messageInput.value.trim();
        if (!content) return;

        try {
            const response = await fetch('/api/messages/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': token,
                },
                body: JSON.stringify({ chat_id: chatID, content }),
            });
            const data = await response.json();
            if (response.ok) {
                messageInput.value = '';
                // Сообщение будет получено через WebSocket
            } else {
                chatError.innerText = data.error || 'Ошибка отправки сообщения';
            }
        } catch (error) {
            chatError.innerText = 'Ошибка сети';
        }
    });
});