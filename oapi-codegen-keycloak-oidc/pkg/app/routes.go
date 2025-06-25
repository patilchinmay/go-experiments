package app

import (
	"context"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	oapiechomiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/patilchinmay/go-experiments/oapi-codegen-keycloak-oidc/api/spec"
	slogecho "github.com/samber/slog-echo"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (app *application) routes() http.Handler {
	// Set up a basic Echo router
	e := echo.New()

	// We need spec in order to validate our requests
	swagger, err := spec.GetSwagger()
	if err != nil {
		app.logger.Error("Failed to get openapi spec: " + err.Error())
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Middlewares
	e.Use(slogecho.New(app.logger))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.BodyLimit("2M"))

	openapiroutes := e.Group("")

	// validation middleware to check all requests against the OpenAPI schema.
	validator := oapiechomiddleware.OapiRequestValidatorWithOptions(swagger, &oapiechomiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
				echoCtx := oapiechomiddleware.GetEchoContext(c)

				// Authenticate and authorize token and populate context with user info, scopes, roles and groups
				err := app.AuthenticateAndAuthorize(echoCtx, input.Scopes)
				if err != nil {
					return err
				}

				return nil
			},
		},
	})

	// Middleware to verify token and populate context
	openapiroutes.Use(validator)

	// Implementation of our openapi route handlers
	handler := app.handler

	// Register our openapi route implementations with the echo router
	spec.RegisterHandlers(openapiroutes, handler)

	// Serve the OpenAPI spec file
	e.GET("/api/spec/spec.yaml", func(c echo.Context) error {
		return c.File("./api/spec/spec.yaml")
	})

	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URLs = []string{"/api/spec/spec.yaml"} // Point to your OpenAPI spec
	}))

	return e
}
