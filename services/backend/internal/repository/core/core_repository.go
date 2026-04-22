package repository

import (
	"errors"
	"fmt"
	"time"

	"drivee/internal/domain"
	. "drivee/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ReadConfig(clientIP string) (*Config, error) {
	db, err := gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("Не удалось подключиться к базе:" + err.Error())
	}

	var config Config
	result := db.Where("ip = ?", clientIP).Limit(1).Find(&config)
	if result.RowsAffected == 0 {
		config = Config{
			Login:     "defaultUser",
			Password:  "defaultPass",
			TypeDB:    "sqlite",
			PathDB:    "chat.db",
			IP:        clientIP,
			CreatedAt: time.Now(),
		}

		// Сохраняем новую конфигурацию в базу данных
		if err := db.Create(&config).Error; err != nil {
			return nil, errors.New("Не удалось сохранить дефолтную конфигурацию: " + err.Error())
		}
	}

	return &config, nil
}

// UpdateConfig - функция для обновления конфигурации в базе данных
func UpdateConfig(config *Config) error {
	db, err := gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	if err != nil {
		return errors.New("не удалось подключиться к базе: " + err.Error())
	}

	// Поиск существующей конфигурации по IP
	var existingConfig Config
	result := db.Where("ip = ?", config.IP).First(&existingConfig)

	// Если запись не найдена, создаем новую, иначе обновляем существующую
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := db.Create(config).Error; err != nil {
			return errors.New("не удалось сохранить конфигурацию: " + err.Error())
		}
	} else if result.Error != nil {
		// Ошибка при выполнении запроса
		return errors.New("ошибка при поиске конфигурации: " + result.Error.Error())
	} else {
		// Обновляем существующую конфигурацию
		existingConfig.Login = config.Login
		existingConfig.Password = config.Password
		existingConfig.TypeDB = config.TypeDB
		existingConfig.PathDB = config.PathDB

		if err := db.Save(&existingConfig).Error; err != nil {
			return errors.New("не удалось обновить конфигурацию: " + err.Error())
		}
	}

	return nil
}

var ErrUserExists = errors.New("user already exists")

func RegisterUser(hub *domain.Hub, username, email, password string) (uint, error) {
	user := User{
		Username: username,
		Email:    email,
	}

	// Хешируем пароль перед сохранением
	if err := user.HashPassword(password); err != nil {
		return 0, err
	}

	// Проверяем, существует ли пользователь с таким именем или email
	var existingUser User
	if err := hub.DB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return 0, ErrUserExists
	}

	// Сохраняем пользователя в базе данных
	if err := hub.DB.Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func CreateCoreDB() (*gorm.DB, error) {
	// Подключаемся к SQLite (файл создастся автоматически)
	db, err := gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("Не удалось подключиться к базе:" + err.Error())
	}

	// Автоматически создаём/обновляем таблицу
	if err := db.AutoMigrate(&Message{}, &Config{}, &User{}); err != nil {
		return nil, errors.New("Ошибка миграции: " + err.Error())
	}

	fmt.Println("Таблицы успешно созданы или обновлены.")

	return db, nil
}
