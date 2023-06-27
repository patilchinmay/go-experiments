package logger

import (
	"os"
	"strconv"

	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func getLogLevel() string {
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
		Logger.Info().Str("logLevel", logLevel).Msg("Loaded logLevel from env var")
	default:
		logLevel = "info"
		Logger.Error().Msg("Invalid LOGLEVEL, setting default logLevel to info")
	}

	return logLevel
}

func getJsonLogs() bool {
	// json logging
	jsonLogs, err := strconv.ParseBool(os.Getenv("JSONLOGS"))
	if err != nil {
		jsonLogs = true
		Logger.Error().Msg("Failed to parse JSONLOGS, setting default jsonLogs to true")
	}

	return jsonLogs
}

func InitiateLogger() zerolog.Logger {
	// Create Logger
	Logger = httplog.NewLogger("Main", httplog.Options{
		JSON: true, LogLevel: "debug",
	})

	// Redefine Logger with proper config
	Logger = httplog.NewLogger("My App", httplog.Options{
		JSON:     getJsonLogs(),
		Concise:  true,
		LogLevel: getLogLevel(),
		Tags: map[string]string{
			"env": os.Getenv("ENV"),
		},
	})

	return Logger
}
