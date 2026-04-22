package websocket

import (
	"log"
	"net/http"

	. "drivee/internal/domain" // твоя доменная папка

	"github.com/gorilla/websocket"
)

// Глобальный Upgrader — объявляем один раз
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Разрешаем подключаться со всех источников (для разработки)
		// В продакшене лучше указывать конкретный origin
		return true
	},
	// Можно добавить лимиты:
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs — основной обработчик WebSocket соединения
func ServeWS(hub *Hub, clientIP string, w http.ResponseWriter, r *http.Request) {
	// Проверяем, что уже есть подключённый клиент
	// if hub.Client != nil {
	// 	http.Error(w, "Чат уже занят другим пользователем", http.StatusConflict)
	// 	log.Println("Попытка подключения второго клиента отклонена")
	// 	return
	// }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда WebSocket:", err)
		return
	}

	client := &Client{
		Hub:  hub,
		Conn: conn,
	}

	hub.Client = client

	log.Println("Клиент успешно подключился к чату с ИИ")
	hub.SendHistoryToClient(clientIP)
	// Запускаем чтение сообщений от пользователя
	go client.ReadPump(clientIP)
}
