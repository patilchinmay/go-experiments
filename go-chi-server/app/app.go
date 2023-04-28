package app

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

type App struct {
	logger zerolog.Logger
}

func New(logger zerolog.Logger) *App {
	return &App{logger: logger}
}

func (a *App) Ping(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())
	oplog.Info().Msg("Pong")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte("Pong"))
}

func (a *App) CreateApp() *chi.Mux {
	router := chi.NewRouter()

	router.Use(httplog.RequestLogger(a.logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/health"))

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}
		// Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
		// ExposedHeaders:   []string{"Link"},
		// AllowCredentials: false,
		// MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/ping", a.Ping)

	return router
}
