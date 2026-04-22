package domain

import (
	. "drivee/internal/models"
	"drivee/internal/repository/sqlite"
	"drivee/internal/usecase"
	"fmt"

	"log"
)

func (c *Client) ReadPump(clientIP string) {
	defer c.Conn.Close()

	msg := &Msg{}

	for {
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Клиент отключился:", err)
			return
		}
		// 1. Показываем сообщение пользователя самому себе
		c.WriteMSG("User", msg.Text, clientIP)
		log.Println(msg.Text)
		// 2. Получаем ответ от ИИ
		log.Println("Отправка в ии")
		aiText, err := usecase.GetAIResponse(msg.Text)
		if err != nil {
			log.Println("Ошибка ии сообщения пользователя:", err)
			continue
		}
		log.Println("исполнение в бд")
		//исполнение в бд
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

		// ответ аи
		c.WriteMSG("AI", "SQL запрос:"+fmt.Sprintf("%v", aiText), clientIP)
		log.Println(aiText)

		c.WriteMSG("AI", "Результаты запроса: "+fmt.Sprintf("%v", results), clientIP)
		log.Println(results)
	}
}
