package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"

	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
	_ "github.com/patilchinmay/go-experiments/go-chi-server/app/goroutineid"
	_ "github.com/patilchinmay/go-experiments/go-chi-server/app/ping"
	"github.com/patilchinmay/go-experiments/go-chi-server/app/user"
	_ "github.com/patilchinmay/go-experiments/go-chi-server/app/validator"
	"github.com/patilchinmay/go-experiments/go-chi-server/db"
	"github.com/patilchinmay/go-experiments/go-chi-server/server"
	globallogger "github.com/patilchinmay/go-experiments/go-chi-server/utils/logger"
	"github.com/rs/zerolog"
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

	// Redefine Logger with proper config
	logger = globallogger.InitiateLogger()

	// Set resources
	// In latest golang version, GOMAXPROCS is automatically set to runtime.NumCPU
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	logger.Info().Int("numCPU", numCPU).Msg("Available resources")

	// Initialize Database (for dependency injection)
	Db := initializeDB(logger)

	// Create app with routes handlers (uses builder pattern)
	app := app.GetOrCreate().SetupDB(Db.DB).WithLogger(logger).SetupCORS().SetupMiddlewares().SetupNotFoundHandler()

	// Create and setup user subrouter
	user.SetupSubrouter(app.DB, logger)

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

func initializeDB(logger zerolog.Logger) *db.Database {
	Db := db.New(logger)

	return Db
}
