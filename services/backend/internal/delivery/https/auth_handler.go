package https

import (
	"drivee/internal/domain"
	. "drivee/internal/models"
	core "drivee/internal/repository/core"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

func Login(tokenAuth *jwtauth.JWTAuth) http.HandlerFunc { // Добавь storage, если нужно проверять по БД
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		fmt.Println(req.Password, req.Username)

		// Проверка credentials (hardcoded для примера; в реальности — из БД с хэшированием паролей, например bcrypt)
		// if req.Username != "admin" || req.Password != "password" {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	render.JSON(w, r, Error("Internal error"))
		// 	return
		// }

		exp := time.Now().Add(24 * time.Hour)

		// Создай claims с ролью
		claims := map[string]interface{}{
			"username": req.Username,
			"exp":      exp,
		}

		_, tokenString, err := tokenAuth.Encode(claims)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Error("Internal error"))
			return
		}

		// Возвращаем токен
		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Expires:  exp,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Path:     "/",
		})
		// json.NewEncoder(w).Encode(map[string]string{
		// 	"token": tokenString,
		// })
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, OK())
	}
}

func Register(hub *domain.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Базовая валидация
		if req.Username == "" || req.Email == "" || req.Password == "" { // Исправлено условие
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Создание пользователя
		userID, err := core.RegisterUser(hub, req.Username, req.Email, req.Password) // Вызов функции регистрации
		if err != nil {
			if err.Error() == "user already exists" {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"message": "User already exists",
					"user_id": userID,
				})
			} else {
				log.Println(err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Успешная регистрация
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User registered successfully",
			"user_id": userID,
		})
	}
}
