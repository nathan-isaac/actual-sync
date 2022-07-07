package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (it *RouteHandler) GetMode(c echo.Context) error {
	return c.String(http.StatusOK, it.Config.ModeString())
}
