package main

import (
	"log"
	"os"

	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
	"github.com/patilchinmay/go-experiments/go-chi-server/server"
)

func main() {
	// If we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create logger
	logger := httplog.NewLogger("My App")

	// Load env variables.
	err := godotenv.Load()
	if err != nil {
		logger.Error().Msg(".env file is not found")
	}

	// Set app env
	env := os.Getenv("ENV")
	logger.Info().Str("ENV", env).Msg("")

	// Logger level setting
	var jsonLogs bool
	var logLevel string
	if env == "local" {
		jsonLogs = false
		logLevel = "debug"
	} else {
		jsonLogs = true
		logLevel = "info"
	}

	// Redefine Logger with proper cconfig
	logger = httplog.NewLogger("My App", httplog.Options{
		JSON:     jsonLogs,
		Concise:  true,
		LogLevel: logLevel,
		Tags: map[string]string{
			"env": env,
		},
	})

	// Start server
	server.New(logger).Serve()
}
