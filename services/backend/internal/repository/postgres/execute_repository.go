package postgres

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

// DB - это глобальная переменная для подключения к базе данных
var DB *sql.DB

// InitializeDB инициализирует подключение к базе данных
func InitializeDB(dataSourceName string) error {
    var err error
    DB, err = sql.Open("postgres", dataSourceName)
    if err != nil {
        return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
    }

    // Проверяем соединение
    if err = DB.Ping(); err != nil {
        return fmt.Errorf("не удалось установить соединение: %w", err)
    }

    return nil
}

// ExecuteSQL выполняет SQL-запрос
func ExecuteSQL(aiText string) error {
    _, err := DB.Exec(aiText)
    if err != nil {
        return fmt.Errorf("ошибка при выполнении SQL-запроса: %w", err)
    }
    return nil
}
