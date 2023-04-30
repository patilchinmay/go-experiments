package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/patilchinmay/go-experiments/go-chi-server/app"
)

// Why struct?
// It is the not necessary to use struct, but it helps to use it.
// Since it helps encapsulate the methods required for server's operation.
// Why accept logger?
// Dependency inversion. It is easier to pass a main logger from the main function.
type Server struct {
	logger zerolog.Logger
	server http.Server
}

// Why create constructor?
// https://web3.coach/golang-why-you-should-use-constructors
// Make refactoring much easier down the line. We can add sensible defaults here.
// We can add validation checks and other necessary logic here down the line.
func New(
	logger zerolog.Logger,
) *Server {
	return &Server{logger: logger,
		server: http.Server{},
	}
}

func (s *Server) Serve() {
	server := s.server

	port := s.getPort()
	host := s.getHost()

	// TODO: We can move these into separate functions (and use a builder paattern)
	server.Addr = host + ":" + port
	server.ReadTimeout = 5 * time.Second
	server.Handler = app.New(s.logger).CreateApp()

	s.logger.Info().Str("Addr", server.Addr).Msg("Listening")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// Error starting or closing listener
		// using Fatal makes sure that the program exits with a status code os.Exit(1) (e.g. when the port is already in use)
		// this helps docker/k8s know that the program is unhealthy and it can take further actions such as restarting the container
		// e.g. when a port is in use, we would like the program to exit fast rather than existing without doing anything
		s.logger.Fatal().Err(err).Msg("Failed to listen and serve")
	} else {
		s.logger.Info().Msg("Server stopped listening")
	}
}

func (s *Server) Shutdown() {
	// Timeout for graceful shutdown
	// Why do we need a timeout context?
	// server.Shutdown does not interrupt active coonections.
	// It works by first closing all open listeners, then closing all idle connections,
	// and then **waiting indefinitely** for active connections to return to idle and then shut down.
	// If the provided context expires before the shutdown is complete,
	// Shutdown returns the context’s error,
	// otherwise it returns any error returned from closing the Server’s underlying Listener(s).
	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown
	if err := s.server.Shutdown(shutdownTimeoutCtx); err != nil {
		s.logger.Fatal().Err(err).Msg("Server shutdown failed")
	}
	s.logger.Info().Msg("Server exited properly")
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
