package http

import (
	// . "drivee/internal/models"
	"drivee/pkg/ip"
	. "drivee/internal/models"

	repository "drivee/internal/repository/core"
	"encoding/json"
	"net/http"
)

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	clientIP := ip.GetRealIP(r)

	config, err := repository.ReadConfig(clientIP)
	if err != nil {
		return
	}

	// Кодируем конфигурацию в JSON и отправляем в ответ
	err = json.NewEncoder(w).Encode(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок ответа
    w.Header().Set("Content-Type", "application/json")

    // Декодируем JSON из тела запроса
    var config Config
    if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
        http.Error(w, "Неверный формат данных", http.StatusBadRequest)
        return
    }

    // Обновляем конфигурацию
    if err := repository.UpdateConfig(&config); err != nil {
        http.Error(w, "Ошибка при обновлении конфигурации: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Отправляем успешный ответ
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Конфигурация обновлена успешно"})

}

