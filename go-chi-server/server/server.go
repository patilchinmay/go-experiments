package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
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
	// Create context that listens for the interrupt signal from the OS.
	// We do this before everything else as
	// We should be capable of catching signal as soon as the program starts
	// https://henvic.dev/posts/signal-notify-context/
	// https://millhouse.dev/posts/graceful-shutdowns-in-golang-with-signal-notify-context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

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
	// The server is started on a separate goroutine
	// ListenAndServe is a blocking function,
	// It fields all incoming requests on separate goroutine
	// https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	go func() {
		logger.Info().Msg("Listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Error starting or closing listener
			logger.Error().Err(err).Msg("Failed to listen and serve")
		} else {
			// This block is only executed when the blocking function
			// ListenAndServe stops and the server is close (ErrServerClosed)
			logger.Info().Msg("Server shutdown successfully")
		}
	}()

	// Block main goroutine to wait for receiving interrupt signal
	<-ctx.Done()

	// Received interrupt signal. Proceed to graceful shutdown
	// stop()
	logger.Info().Msg("Shutting down gracefully")

	// Timeout for graceful shutdown
	// Why need a timeout context?
	// srv.Shutdown does not interrupt active coonections.
	// It works by first closing all open listeners, then closing all idle connections,
	// and then **waiting indefinitely** for active connections to return to idle and then shut down.
	// If the provided context expires before the shutdown is complete,
	// Shutdown returns the context’s error,
	// otherwise it returns any error returned from closing the Server’s underlying Listener(s).
	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown
	if err := server.Shutdown(shutdownTimeoutCtx); err != nil {
		logger.Fatal().Err(err).Msg("Server shutdown failed")
	}
	logger.Info().Msg("Server exited properly")
}

func (s *Server) Shutdown() {

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
