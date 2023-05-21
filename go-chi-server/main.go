package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"strconv"

	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
	_ "github.com/patilchinmay/go-experiments/go-chi-server/app/goroutineid"
	_ "github.com/patilchinmay/go-experiments/go-chi-server/app/ping"
	_ "github.com/patilchinmay/go-experiments/go-chi-server/app/users"
	_ "github.com/patilchinmay/go-experiments/go-chi-server/app/validator"
	"github.com/patilchinmay/go-experiments/go-chi-server/server"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	// We do this at the beginning of the program as
	// We should be capable of catching signal as soon as the program starts
	// https://henvic.dev/posts/signal-notify-context/
	// https://millhouse.dev/posts/graceful-shutdowns-in-golang-with-signal-notify-context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create logger
	logger := httplog.NewLogger("Main", httplog.Options{
		JSON: true, LogLevel: "debug",
	})

	// Load env variables.
	err := godotenv.Load()
	if err != nil {
		logger.Error().Msg(".env file is not found")
	}

	// Set app env
	env := os.Getenv("ENV")
	logger.Info().Str("ENV", env).Msg("")

	// json logging
	jsonLogs, err := strconv.ParseBool(os.Getenv("JSONLOGS"))
	if err != nil {
		jsonLogs = true
		logger.Error().Msg("Failed to parse JSONLOGS, setting default jsonLogs to true")
	}

	// log level setting and validation
	logLevel := os.Getenv("LOGLEVEL")
	switch logLevel {
	case
		"trace",
		"debug",
		"info",
		"warn",
		"error",
		"critical":
		logger.Info().Str("logLevel", logLevel).Msg("Loaded logLevel from env var")
	default:
		logLevel = "info"
		logger.Error().Msg("Invalid LOGLEVEL, setting default logLevel to info")
	}

	// Redefine Logger with proper config
	logger = httplog.NewLogger("My App", httplog.Options{
		JSON:     jsonLogs,
		Concise:  true,
		LogLevel: logLevel,
		Tags: map[string]string{
			"env": env,
		},
	})

	// Set resources
	// In latest golang version, GOMAXPROCS is automatically set to runtime.NumCPU
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	logger.Info().Int("numCPU", numCPU).Msg("Available resources")

	// Create app with routes handlers (uses builder pattern)
	app := app.GetOrCreate().WithLogger(logger).SetupCORS().SetupMiddlewares().SetupNotFoundHandler()

	// Mounts subrouters on main app/router
	app.MountSubrouters()

	// Create server
	server := server.New().WithLogger(logger).WithHandlers(app.Router)

	// The server is started on a separate goroutine as
	// ListenAndServe is a blocking function,
	// It fields all incoming requests on separate goroutine
	// https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	go server.Serve()

	// Block main goroutine to wait for receiving os interrupt signal
	<-ctx.Done()

	// Received interrupt signal. Proceed to graceful shutdown
	logger.Info().Msg("Received interrupt, shutting down gracefully")
	server.Shutdown()
}
