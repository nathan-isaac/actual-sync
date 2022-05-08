package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (it *RouteHandler) GetMode(c echo.Context) error {
	return c.String(http.StatusOK, it.Config.ModeString())
}
