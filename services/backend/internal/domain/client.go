package domain

import (
	. "drivee/internal/models"
	"drivee/internal/repository/sqlite"
	"drivee/internal/usecase"
	"fmt"

	"log"
	"time"
)

var msg struct {
	Text string `json:"text"`
	// Type string `json:"type"`
}

func (c *Client) ReadPump(clientIP string) {
	defer c.Conn.Close()

	for {
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Клиент отключился:", err)
			return
		}

		// 1. Показываем сообщение пользователя самому себе
		userMsg := Message{
			Username:  "User",
			Text:      msg.Text,
			IP:        clientIP,
			CreatedAt: time.Now(),
		}
		if err := c.Hub.DB.Create(&userMsg).Error; err != nil {
			log.Println("Ошибка сохранения сообщения пользователя:", err)
		}

		c.Conn.WriteJSON(userMsg)

		// 2. Получаем ответ от ИИ
		log.Println("Отправка в ии")
		aiText, err := usecase.GetAIResponse(msg.Text)
		if err != nil {
			log.Println("Ошибка ии сообщения пользователя:", err)
			continue
		}

		db, err := sqlite.InitializeDB("chat.db")
		if err != nil {
			log.Println("Ошибка ии сообщения пользователя:", err)
		}

		rows, err := db.ExecuteSQL(aiText)
		if err != nil {
			log.Println("Ошибка ии сообщения пользователя:", err)
		}
		if err != nil {
			// aiText = err.Error()
			// c.Conn.WriteJSON("неа")

			// ошибка
			continue
		}
		defer rows.Close()

		var results []map[string]interface{} // Используем срез для хранения результатов

		columns, _ := rows.Columns() // Получаем названия колонок
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for rows.Next() {
			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				log.Println("Ошибка сканирования строки:", err)
				continue
			}

			result := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				result[col] = val
			}

			results = append(results, result) // Добавляем результат в срез
		}

		if err := rows.Err(); err != nil {
			log.Println("Ошибка при итерации по строкам:", err)
			continue
		}

		aiMsg := Message{
			Username:  "AI",
			Text:      "SQL запрос:" + fmt.Sprintf("%v", msg.Text),
			IP:        clientIP,
			CreatedAt: time.Now(),
		}

		if err := c.Hub.DB.Create(&aiMsg).Error; err != nil {
			log.Println("Ошибка сохранения ответа ИИ:", err)
		}
		
		c.Conn.WriteJSON(aiMsg)

		aiMsg = Message{
			Username:  "AI",
			Text:      "Результаты запроса: " + fmt.Sprintf("%v", results),
			IP:        clientIP,
			CreatedAt: time.Now(),
		}

		if err := c.Hub.DB.Create(&aiMsg).Error; err != nil {
			log.Println("Ошибка сохранения ответа ИИ:", err)
		}

		// 3. Отправляем ответ ИИ пользователю
		c.Conn.WriteJSON(aiMsg)
	}
}
