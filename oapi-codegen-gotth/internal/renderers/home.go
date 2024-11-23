package renderers

import (
	"github.com/labstack/echo/v4"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/internal/views"
)

func HomePage(c echo.Context) error {
	component := views.HomePage()
	return component.Render(c.Request().Context(), c.Response().Writer)
}
