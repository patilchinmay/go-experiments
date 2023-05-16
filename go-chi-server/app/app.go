package app

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

type App struct {
	logger     zerolog.Logger
	Router     *chi.Mux
	Subrouters []Subrouter
}

var app *App

// GetOrCreate returns a pointer to App using singleton pattern.
// If app exists, it returns it. If not, it creates it and returns it
func GetOrCreate() *App {
	if app == nil {
		app = &App{
			Router: chi.NewRouter(),
		}
	}
	return app
}

// Discard will remove the reference to app so that it can be garbage collected. In other words, it deletes the singleton instance of *App.
func Discard() {
	if app != nil {
		app = nil
	}
}

func (a *App) WithLogger(logger zerolog.Logger) *App {
	a.logger = logger
	return a
}

func (a *App) SetupMiddlewares() *App {
	// httplog.RequestLogger sets up RequestId and Recoverer as well
	a.Router.Use(httplog.RequestLogger(a.logger))
	a.Router.Use(middleware.Heartbeat("/health"))

	return a
}

func (a *App) SetupCORS() *App {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	a.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}
		// Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
		// ExposedHeaders:   []string{"Link"},
		// AllowCredentials: false,
		// MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	return a
}

func (a *App) SetupNotFoundHandler() *App {
	// https://github.com/go-chi/chi/issues/780
	a.Router.HandleFunc("/", a.Router.NotFoundHandler())

	return a
}
