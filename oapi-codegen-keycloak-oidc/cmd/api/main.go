package main

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/lmittmann/tint"
	"github.com/patilchinmay/go-experiments/oapi-codegen-keycloak-oidc/pkg/app"
	"github.com/patilchinmay/go-experiments/oapi-codegen-keycloak-oidc/pkg/auth"
	"github.com/patilchinmay/go-experiments/oapi-codegen-keycloak-oidc/pkg/handler"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, AddSource: true}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	// Load .env file for local development (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Load Keycloak configuration from environment variables
	var keycloakConfig auth.KeycloakConfig
	err := envconfig.Process("", &keycloakConfig)
	if err != nil {
		return err
	}

	v := auth.NewKeycloakValidator(logger, keycloakConfig)

	// Create Handler
	h := handler.NewPetStore(logger)

	// Load Application configuration from environment variables
	var cfg app.Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		return err
	}

	app := app.New(logger, cfg, h, v)

	return app.ServeHTTP()
}
