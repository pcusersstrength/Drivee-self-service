package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AIResponse описывает структуру ответа от API
type AIResponse struct {
	Success  bool   `json:"success"`
	Question string `json:"question"`
	Dialect  string `json:"dialect"`
	SQL      string `json:"sql"`
}

const (
	apiURL   = "http://api.daemuun.ru:8000/api/sql"
	apiToken = "9Xov01Gwyc1rUPNT86rFIEQ4HUlt1uW88tnSA683MIc=" // Замените на реальный токен
)

// RequestBody описывает структуру тела запроса (если API ожидает JSON)
type RequestBody struct {
	Text string `json:"text"`
}

func GetAIResponse(text string) (string, error) {
	// 1. Подготовка данных для отправки
	requestData := RequestBody{
		Text: text,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 2. Создание HTTP запроса
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Установка заголовков
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken) // Обычно токены передаются так

	// 4. Выполнение запроса с таймаутом
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 5. Проверка статус-кода
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api returned error status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 6. Чтение и парсинг ответа
	var aiResp AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 7. Проверка бизнес-логики успеха
	if !aiResp.Success {
		return "", fmt.Errorf("api returned success=false")
	}

	// Возвращаем SQL запрос как результат
	return aiResp.SQL, nil
}
