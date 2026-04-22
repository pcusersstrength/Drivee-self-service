package models

import (
	"time"
)

// Message — структура сообщения для хранения в БД и отправки по WebSocket
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"` // "You" или "AI"
	Text      string    `json:"text"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

type Msg struct {
	Text string `json:"text"`
	// Type string `json:"type"`
}

type Config struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Login     string    `json:"username"` 
	Password  string    `json:"pass"`
	TypeDB    string    `json:"type_db"`
	PathDB    string    `json:"path_db"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}
