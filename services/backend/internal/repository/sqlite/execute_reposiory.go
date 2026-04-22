package sqlite

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3" // Импортируем драйвер SQLite
)

// DB - это глобальная переменная для подключения к базе данных
type SQLite struct {
    DB *sql.DB
}

// ExecuteSQL выполняет SELECT-запрос и возвращает результат
func (s *SQLite) ExecuteSQL(query string) (*sql.Rows, error) {
    rows, err := s.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("ошибка при выполнении SELECT-запроса: %w", err)
    }
    return rows, nil
}

// InitializeDB инициализирует подключение к базе данных и возвращает экземпляр SQLiteDB
func InitializeDB(dataSourceName string) (*SQLite, error) {
    db, err := sql.Open("sqlite3", dataSourceName)
    if err != nil {
        return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
    }

    // Проверяем соединение
    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("не удалось установить соединение: %w", err)
    }

    return &SQLite{DB: db}, nil
}

// CloseDB закрывает соединение с базой данных
func (s *SQLite) Close() error {
    if s.DB != nil {
        return s.DB.Close()
    }
    return nil
}