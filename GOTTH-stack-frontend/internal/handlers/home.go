package handlers

import (
	"github.com/patilchinmay/go-experiments/gotth-stack-frontend/internal/views"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	component := views.Home()
	return component.Render(c.Request().Context(), c.Response().Writer)
}
