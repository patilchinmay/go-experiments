package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	DefaultHost = "0.0.0.0"
	DefaultPort = "8080"
)

// Why struct?
// It is the not necessary to use struct, but it helps to use it.
// Since it helps encapsulate the methods required for server's operation.
// Why accept logger?
// Dependency inversion. It is easier to pass a main logger from the main function.
type Server struct {
	logger zerolog.Logger
	host   string
	port   string
	server http.Server
}

// Why create constructor?
// https://web3.coach/golang-why-you-should-use-constructors
// Make refactoring much easier down the line. We can add sensible defaults here.
// We can add validation checks and other necessary logic here down the line.
func New() *Server {
	s := &Server{
		server: http.Server{},
	}
	s.host = s.setDefaultHost()
	s.port = s.setDefaultPort()

	return s
}

// WithLogger sets the logger using builder pattern
func (s *Server) WithLogger(logger zerolog.Logger) *Server {
	s.logger = logger
	return s
}

// WithHandlers sets the route handlers using builder pattern
func (s *Server) WithHandlers(app http.Handler) *Server {
	s.server.Handler = app
	return s
}

// WithReadTimeout sets the read timeout for server using builder pattern
func (s *Server) WithReadTimeout(duration time.Duration) *Server {
	s.server.ReadTimeout = duration
	return s
}

// WithHost sets the server host address using builder pattern
func (s *Server) WithHost(host string) *Server {
	s.host = host
	return s
}

// WithPort sets the server port using builder pattern
func (s *Server) WithPort(port string) *Server {
	s.port = port
	return s
}

// Serve servers requests on the mentioned host and port
func (s *Server) Serve() {

	s.server.Addr = s.host + ":" + s.port

	s.logger.Info().Str("Addr", s.server.Addr).Msg("Listening")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// Error starting or closing listener
		// using Fatal makes sure that the program exits with a status code os.Exit(1) (e.g. when the port is already in use)
		// this helps docker/k8s know that the program is unhealthy and it can take further actions such as restarting the container
		// e.g. when a port is in use, we would like the program to exit fast rather than existing without doing anything
		s.logger.Fatal().Err(err).Msg("Failed to listen and serve")
	} else {
		s.logger.Info().Msg("Server stopped listening")
	}
}

// Shutdown shuts down the server with a timeout of 5 seconds
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

// setDefaultHost sets the default host to "0.0.0.0" if HOST env var is not set
func (s *Server) setDefaultHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = DefaultHost
	}
	return host
}

// setDefaultPort sets the default port to "8080" is PORT env var is not set
func (s *Server) setDefaultPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	return port
}
