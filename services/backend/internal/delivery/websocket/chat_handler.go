package websocket

import (
	"log"
	"net/http"

	. "drivee/internal/domain" // твоя доменная папка

	"github.com/gorilla/websocket"
)

// Глобальный Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWS(hub *Hub, clientIP string, w http.ResponseWriter, r *http.Request) {

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

	go client.ReadPump(clientIP)
}
