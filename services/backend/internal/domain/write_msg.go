package domain

import (
	. "drivee/internal/models"
	"log"
	"time"
)

// 1. Показываем сообщение пользователя самому себе

func (c *Client) WriteMSG(user, text, clientIP string) error {
	userMsg := Message{
		Username:  user,
		Text:      text,
		IP:        clientIP,
		CreatedAt: time.Now(),
	}

	if err := c.Hub.DB.Create(&userMsg).Error; err != nil {
		log.Println("Ошибка сохранения сообщения пользователя:", err)
	}

	c.Conn.WriteJSON(userMsg)
	return nil
}
