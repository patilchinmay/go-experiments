//revive:disable:package-comments
package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
	middleware "github.com/oapi-codegen/echo-middleware"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/internal/handlers"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/pkg/spec/generated"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, AddSource: true}))

	slog.SetDefault(logger)

	port := flag.String("port", "3000", "Port for test HTTP server")
	flag.Parse()

	swagger, err := generated.GetSwagger()
	if err != nil {
		slog.Error("error loading swagger spec", "error", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	petStore := handlers.NewPetStore()

	// This is how you set up a basic Echo router
	e := echo.New()

	// Middlewares

	// Request Logger
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		//revive:disable:unused-parameter
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(echomiddleware.Recover())

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	generated.RegisterHandlers(e, petStore)

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(net.JoinHostPort("0.0.0.0", *port)))
}
