package main

import (
	"log"
	"net/http"

	"drivee/pkg/ip"
	. "drivee/internal/delivery/http"
	. "drivee/internal/delivery/websocket"
	. "drivee/internal/domain"
)

func main() {
	hub := NewHub()

	http.HandleFunc("/api/get_config", GetConfigHandler)

	http.HandleFunc("/api/update_config", UpdateConfigHandler)

	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Разрешаем только нужный домен (или несколько)
		allowedOrigins := map[string]bool{
			"https://higu.su":     true,
			"http://higu.su":      true,
			"https://www.higu.su": true,
		}

		if origin != "" && allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if origin == "" {
			// Для теста можно разрешить всё, но в продакшене — не рекомендуется
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Upgrade, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Protocol, Sec-WebSocket-Extensions")
		// w.Header().Set("Access-Control-Allow-Credentials", "true") // если используешь cookies/auth

		// Обработка preflight OPTIONS (очень важно для Safari!)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		clientIP := ip.GetRealIP(r)
		log.Printf("WebSocket connection from: %s (RemoteAddr: %s)", clientIP, r.RemoteAddr)

		ServeWS(hub, clientIP, w, r)
	})

	// Простая страница для теста чата
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
