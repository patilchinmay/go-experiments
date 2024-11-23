//revive:disable:package-comments
package main

import (
	"context"
	"flag"
	"io/fs"
	"log/slog"
	"net"
	"os"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
	middleware "github.com/oapi-codegen/echo-middleware"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/internal/handlers"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/internal/renderers"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/pkg/spec/generated"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/public"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, AddSource: true}))

	slog.SetDefault(logger)

	port := flag.String("port", "3000", "Port for test HTTP server")
	flag.Parse()

	// https://github.com/oapi-codegen/oapi-codegen/issues/514
	openapi3filter.RegisterBodyDecoder("text/yaml", openapi3filter.FileBodyDecoder)

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

	// Global middlewares

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
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())

	// Static files - no validation needed
	e.StaticFS("/public", fs.FS(public.Assets))

	// Web UI routes - no validation needed
	webUI := e.Group("/ui")
	// Home page route for Web UI clients
	webUI.GET("", renderers.HomePage)

	// API routes - with OpenAPI validation
	api := e.Group("")
	// Use our validation middleware to check all requests against the OpenAPI schema.
	api.Use(middleware.OapiRequestValidator(swagger))
	// We now register our petStore above as the handler for the interface
	generated.RegisterHandlers(api, petStore)

	// Print all routes for debugging
	for _, route := range e.Routes() {
		slog.Info("registered route",
			"method", route.Method,
			"path", route.Path,
		)
	}

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(net.JoinHostPort("127.0.0.1", *port)))
}
