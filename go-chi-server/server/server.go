package server

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/patilchinmay/go-experiments/go-chi-server/app"
)

// Why struct?
//
// Why accept logger?
type Server struct {
	logger zerolog.Logger
}

// Why create constructor?
func New(
	logger zerolog.Logger,
) *Server {
	return &Server{logger: logger}
}

func (s *Server) Serve() {
	logger := s.logger
	port := s.getPort()
	host := s.getHost()

	server := http.Server{
		Addr:        host + ":" + port,
		ReadTimeout: 5 * time.Second,
		Handler:     app.New(logger).CreateApp(),
	}
	logger = logger.With().Str("Addr", server.Addr).Logger()

	// Starting server
	logger.Info().Msg("Listening")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error().Err(err).Msg("Failed to listen and serve")
	} else {
		logger.Info().Msg("Server shutdown successfully")
	}
	// Handle interrupt signals
}

func (s *Server) getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	s.logger.Info().Str("PORT", port).Msg("Deploying webhook")

	return port
}

func (s *Server) getHost() string {
	host := os.Getenv("HOST")

	if host == "" {
		host = "0.0.0.0"
	}
	s.logger.Info().Str("HOST", host).Msg("Deploying webhook")

	return host
}
