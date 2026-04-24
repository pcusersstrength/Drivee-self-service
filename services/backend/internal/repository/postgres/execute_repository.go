package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // Современный драйвер для PostgreSQL
)

// Postgres - структура для управления подключением
type Postgres struct {
	DB *sql.DB
}

// ExecSQL выполняет SELECT-запрос и возвращает результат
func (p *Postgres) Exec(query string) (string, error) {
	rows, err := p.DB.Query(query)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении SELECT-запроса: %w", err)
	}

	defer rows.Close()

	columns, _ := rows.Columns() // Получаем названия колонок
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	var results []map[string]interface{} // Используем срез для хранения результатов

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
		return "", err
	}

	return fmt.Sprintf("%v", results), nil

}

func (p *Postgres) ExecuteSQL(query string) (*sql.Rows, error) {
	// if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "SELECT") {
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении SELECT-запроса: %w", err)
	}
	return rows, nil
	// } else {
	// 	result, err := p.DB.Exec(query)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	n, err := result.RowsAffected()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	log.Println(n)
	// 	return nil, nil

	// }
	// return nil,nil
}

// InitializeDB инициализирует подключение к PostgreSQL
// dataSourceName должен быть в формате:
// "postgres://myuser:mydatabase@localhost:5432/mydatabase?sslmode=disable"
func InitializeDB(dataSourceName string) (*Postgres, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть драйвер: %w", err)
	}

	// Проверяем реальное соединение с сервером
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось установить соединение с PostgreSQL: %w", err)
	}

	// Настройки пула соединений (рекомендуется для Postgres)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return &Postgres{DB: db}, nil
}

// Close закрывает соединение с базой данных
func (p *Postgres) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}
