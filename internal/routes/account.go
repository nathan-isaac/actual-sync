package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type NeedsBootstrapData struct {
	Bootstrapped bool `json:"bootstrapped"`
}

type NeedsBootstrapResponse struct {
	Status string             `json:"status"`
	Data   NeedsBootstrapData `json:"data"`
}

func (it *RouteHandler) NeedsBootstrap(c echo.Context) error {
	count, err := it.PasswordStore.Count()

	if err != nil {
		return err
	}

	r := &NeedsBootstrapResponse{
		Status: "ok",
		Data: NeedsBootstrapData{
			Bootstrapped: count > 0,
		},
	}

	return c.JSON(http.StatusOK, r)
}
