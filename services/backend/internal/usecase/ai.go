package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Text    string `json:"q"`
	Dialect string `json:"dialect"`
}

func GetAIResponse(text string) (string, error) {
	// 1. Формирование URL с Query-параметрами
	// API требует параметр "q" в строке запроса (после знака ?)
	u, err := url.Parse(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse api URL: %w", err)
	}

	params := url.Values{}
	params.Add("q", text) // Добавляем текст запроса в query
	// Если API требует dialect тоже в query, раскомментируйте строку ниже:
	// params.Add("dialect", "postgresql")

	u.RawQuery = params.Encode()

	// 2. Создание HTTP запроса (используем GET, так как параметры в URL)
	// Тело (body) для GET запроса обычно не передается
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Установка заголовков
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Accept", "application/json")

	// 4. Выполнение запроса с таймаутом
	client := &http.Client{Timeout: 60 * time.Second}
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

	return aiResp.SQL, nil
}
