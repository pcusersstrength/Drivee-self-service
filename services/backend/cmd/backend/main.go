package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"drivee/internal/delivery/https"
	. "drivee/internal/delivery/websocket"
	. "drivee/internal/domain"
	repository "drivee/internal/repository/core"
	"drivee/pkg/ip"
	sl "drivee/pkg/slogpretty"

	mwLogger "drivee/internal/delivery/middleware/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	log := setupLogger("local")

	db, err := repository.CreateCoreDB()
	if err != nil {
		panic(err)
	}
	hub := NewHub(db)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(mwLogger.New(log))
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://0.0.0.0:*", "https://higu.su", "https://app.higu.su", "drivee-ai.ru", "app.drivee-ai.ru"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Get("/api/get_config", https.GetConfigHandler)
		r.Post("/api/update_config", https.UpdateConfigHandler)
	})

	// WebSocket endpoint
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientIP := ip.GetRealIP(r)
		// log.Info("WebSocket connection", slog.String("Remote Address", r.RemoteAddr), slog.String("client IP", clientIP))

		ServeWS(hub, clientIP, w, r)
	})

	r.Post("/api/get_chart", https.GetChart())

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("public"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	r.Get("/cyrillic-font.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "cyrillic-font.js")
	})
	

	log.Info("starting server", slog.String("address", "localhost:8080"))

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("HTTP server error", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := sl.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
