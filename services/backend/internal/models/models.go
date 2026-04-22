package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
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
	Type string `json:"type"`
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

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// Метод для хеширования пароля (можно использовать bcrypt)
func (u *User) HashPassword(password string) error {
	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Метод для проверки пароля
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
