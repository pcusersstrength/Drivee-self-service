package main

import (
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"drivee/internal/config"
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
	"github.com/go-chi/jwtauth/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var frontendFS embed.FS

func main() {
	cfg := config.MustLoad()

	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	log := setupLogger(cfg.Env)

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
	r.Use(middleware.URLFormat)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://0.0.0.0:*", "https://higu.su", "https://app.higu.su"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/api/get_config", https.GetConfigHandler)
		r.Post("/api/update_config", https.UpdateConfigHandler)
	})
	// WebSocket endpoint
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientIP := ip.GetRealIP(r)
		log.Info("WebSocket connection", slog.String("Remote Address", r.RemoteAddr), slog.String("client IP", clientIP))

		authHeader := r.Header.Get("Authorization")
		log.Info("Authorization header:", slog.String("header", authHeader))

		ServeWS(hub, clientIP, w, r)
	})

	r.Post("/api/auth/login", https.Login(tokenAuth))
	r.Post("/api/auth/register", https.Register(hub))

	// Простая страница для теста чата
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
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
