package domain

import (
	. "drivee/internal/models"
	"drivee/internal/repository/postgres"
	"fmt"

	"log"
)

func (c *Client) ReadPump(clientIP string) {
	defer c.Conn.Close()

	
	for {
		msg := &Msg{}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Клиент отключился:", err)
			return
		}
		if msg.Type == "ping" {
			c.WriteMSG("pong", "system", "pong", clientIP)
			// continue
		} else {
			// 1. Показываем сообщение пользователя самому себе
			c.WriteMSG("", "User", msg.Text, clientIP)
			log.Println(msg.Text)
	
			c.WriteMSG("step1", "AI", "Запрос создан", clientIP)
	
			// 2. Получаем схему для ИИ
			db, err := postgres.InitializeDB("postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable")
			if err != nil {
				log.Println("Ошибка подключения к бд:", err)
			}
	
			schema, err := db.Exec(`SELECT column_name, data_type
					FROM information_schema.columns
					WHERE table_name='anonymized_incity_orders';`)
			if err != nil {
				log.Println("Ошибка выполнения запроса:", err)
				return
			}
	
			// 3. Получаем ответ от ИИ
			log.Println("Отправка в ии")
			aiText, err := GetAIResponseWithSchema(c, msg.Text, schema, clientIP)
			if err != nil {
				log.Println("Ошибка ии сообщения пользователя:", err)
				continue
			}
	
			// 4. Добавляем sql в сообщение
			c.WriteMSG("step2", "AI", fmt.Sprintf("%v", aiText), clientIP)
			log.Println(aiText)
	
			// 5. Исполнение в бд
			log.Println("исполнение в бд")
			rows, err := db.ExecuteSQL(aiText)
			if err != nil {
				log.Println("Ошибка ии сообщения пользователя:", err)
				continue
			}
			if rows != nil {
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
	
				c.WriteMSG("step3", "AI", fmt.Sprintf("%v", results), clientIP)
				log.Println(results)
			} else {
				c.WriteMSG("table", "AI", "Результаты запроса: злодеяние", clientIP)
				log.Println("злодеяние")
			}
		}
	}
}
