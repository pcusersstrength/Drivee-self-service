package domain

import (
	"log"

	"gorm.io/gorm"

	. "drivee/internal/models"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Client *Client
	DB     *gorm.DB
}

type Client struct {
	Conn *websocket.Conn
	Hub  *Hub
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		DB: db,
	}
}

// sendHistoryToClient отправляет последние 50 сообщений владельцу ip
func (h *Hub) SendHistoryToClient(clientIP string) {

	var messages []Message
	h.DB.Where("ip = ?", clientIP).
		Order("created_at ASC").
		Limit(50).
		Find(&messages)

	for _, msg := range messages {
		if err := h.Client.Conn.WriteJSON(msg); err != nil {
			log.Println("Ошибка отправки истории:", err)
			return
		}
	}
}
