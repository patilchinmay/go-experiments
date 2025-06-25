package app

import (
	"log/slog"
	"sync"

	"github.com/patilchinmay/go-experiments/oapi-codegen-keycloak-oidc/pkg/auth"
	"github.com/patilchinmay/go-experiments/oapi-codegen-keycloak-oidc/pkg/handler"
)

type Config struct {
	Host     string `default:"127.0.0.1" envconfig:"HOST"`
	HttpPort int    `default:"9000"      envconfig:"HTTP_PORT"`
}

type application struct {
	config    Config
	logger    *slog.Logger
	handler   handler.HandlerIntf
	validator auth.ValidatorIntf
	wg        sync.WaitGroup
}

var app *application

// New uses the singleton pattern to create/return and instance of application
// It should be called to initialize the application from the main.
func New(
	logger *slog.Logger,
	config Config,
	handler handler.HandlerIntf,
	validator auth.ValidatorIntf,
) *application {
	if app == nil {
		app = &application{
			logger:    logger,
			config:    config,
			handler:   handler,
			validator: validator,
		}
	}

	logger.Info("application initialized")

	return app
}
