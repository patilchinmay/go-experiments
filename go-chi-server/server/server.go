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
	// Create context that listens for the interrupt signal from the OS.
	// We do this before everything else as
	// We should be capable of catching signal as soon as the program starts
	// https://henvic.dev/posts/signal-notify-context/
	// https://millhouse.dev/posts/graceful-shutdowns-in-golang-with-signal-notify-context
	// TODO: Move interrupt handling to main
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := s.logger
	server := s.server

	port := s.getPort()
	host := s.getHost()

	// TODO: We can move these into separate functions (and create a builder paattern)
	server.Addr = host + ":" + port
	server.ReadTimeout = 5 * time.Second
	server.Handler = app.New(logger).CreateApp()

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
			// using Fatal makes sure that the program exits with a status code os.Exit(1) (e.g. when the port is already in use)
			// this helps docker/k8s know that the program is unhealthy and it can take further actions such as restarting the container
			// e.g. when a port is in use, we would like the program to exit fast rather than existing without doing anything
			logger.Fatal().Err(err).Msg("Failed to listen and serve")
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
	// TODO: Move shutdown into its own function
	if err := server.Shutdown(shutdownTimeoutCtx); err != nil {
		logger.Fatal().Err(err).Msg("Server shutdown failed")
	}
	logger.Info().Msg("Server exited properly")
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
