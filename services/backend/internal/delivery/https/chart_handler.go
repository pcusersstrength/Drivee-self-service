package https

import (
	"drivee/internal/repository/postgres"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetChart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Закрываем тело после чтения
		_ = r.Body.Close()

		// Попробуем распарсить как JSON: {"text": "..."}
		if len(body) > 0 {
			var payload struct {
				Sql string `json:"sql"`
			}
			if err := json.Unmarshal(body, &payload); err == nil && payload.Sql != "" {
				// 2. Получаем схему для ИИ
				db, err := postgres.InitializeDB("postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable")
				if err != nil {
					log.Println("Ошибка подключения к бд:", err)
				}

				rows, err := db.ExecuteSQL(payload.Sql)
				if err != nil {
					log.Println("Ошибка подключения к бд:", err)
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
					log.Println(results)
					json.NewEncoder(w).Encode(map[string]any{"chart": results})
					return
				}
				log.Println("Ошибка Unmarshal", err)

			}

			// Здесь можно вернуть конфигурацию (как в вашем оригинальном коде),
			// если текст в body не передан. Пример ниже просто возвращает ошибку:
			http.Error(w, "text not provided in body", http.StatusBadRequest)
		}
	}
}
