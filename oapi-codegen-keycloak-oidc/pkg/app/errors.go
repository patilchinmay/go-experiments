package app

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

// nolint: unused
func (app *application) reportServerError(ctx echo.Context, err error, code int, message string) error {
	app.logger.Error("Server Error",
		slog.String("requestID", ctx.Response().Header().Get(echo.HeaderXRequestID)),
		slog.String("details", err.Error()),
		slog.String("message", message),
	)

	return echo.NewHTTPError(code, message)
}
