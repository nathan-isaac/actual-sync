package routes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type FailureResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type NeedsBootstrapData struct {
	Bootstrapped bool `json:"bootstrapped"`
}

type NeedsBootstrapResponse struct {
	SuccessResponse
	Data NeedsBootstrapData `json:"data"`
}

func (it *RouteHandler) NeedsBootstrap(c echo.Context) error {
	count, err := it.PasswordStore.Count()

	if err != nil {
		return err
	}

	r := &NeedsBootstrapResponse{
		SuccessResponse: SuccessResponse{Status: "ok"},
		Data: NeedsBootstrapData{
			Bootstrapped: count > 0,
		},
	}

	return c.JSON(http.StatusOK, r)
}

type BootstrapRequestBody struct {
	Password core.Password `json:"password"`
}

type BootstrapData struct {
	Token core.Token `json:"token"`
}

type BootstrapResponse struct {
	SuccessResponse
	Data BootstrapData `json:"data"`
}

func (it *RouteHandler) Bootstrap(c echo.Context) error {
	req := new(BootstrapRequestBody)
	if err := c.Bind(req); err != nil {
		return err
	}

	if req.Password == "" {
		r := &FailureResponse{
			Status: "error",
			Reason: "invalid-password",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	count, err := it.PasswordStore.Count()
	if err != nil {
		return err
	}
	if count != 0 {
		r := &FailureResponse{
			Status: "error",
			Reason: "already-bootstrapped",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return err
	}
	err = it.PasswordStore.Add(string(hashed))
	if err != nil {
		return err
	}

	token := uuid.NewString()
	err = it.TokenStore.Add(token)
	if err != nil {
		return err
	}
	r := &BootstrapResponse{
		SuccessResponse: SuccessResponse{Status: "ok"},
		Data:            BootstrapData{Token: token},
	}
	return c.JSON(http.StatusOK, r)
}

type LoginRequestBody struct {
	Password core.Password `json:"password"`
}

type LoginData struct {
	Token interface{} `json:"token"`
}

type LoginSuccessData struct {
	LoginData
	Token core.Token `json:"token"`
}

type LoginSuccessResponse struct {
	SuccessResponse
	Data LoginSuccessData `json:"data"`
}

type LoginFailResponse struct {
	Status string    `json:"status"`
	Data   LoginData `json:"data"`
}

func (it *RouteHandler) Login(c echo.Context) error {
	req := new(LoginRequestBody)
	if err := c.Bind(req); err != nil {
		return err
	}

	hashedPass, err := it.PasswordStore.First()
	if err != nil {
		r := &LoginFailResponse{Status: "ok", Data: LoginData{Token: nil}}
		return c.JSON(http.StatusOK, r)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Password)); err == nil {
		// Right now, tokens are permanent and there's just one in the
		// system. In the future this should probably evolve to be a
		// "session" that times out after a long time or something, and
		// maybe each device has a different token
		token, err := it.TokenStore.First()
		if err != nil {
			return err
		}
		r := &LoginSuccessResponse{
			SuccessResponse: SuccessResponse{Status: "ok"},
			Data:            LoginSuccessData{Token: token},
		}
		return c.JSON(http.StatusOK, r)
	}

	r := &LoginFailResponse{Status: "ok", Data: LoginData{Token: nil}}
	return c.JSON(http.StatusOK, r)
}

type ChangePassRequestBody struct {
	Token    core.Token    `json:"token"`
	Password core.Password `json:"password"`
}

func (it *RouteHandler) ChangePassword(c echo.Context) error {
	req := new(ChangePassRequestBody)
	if err := c.Bind(req); err != nil {
		return err
	}
	val := it.authenticateUser(c, &req.Token)
	if !val {
		r := &FailureResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	if req.Password == "" {
		r := &FailureResponse{
			Status: "error",
			Reason: "invalid-password",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return err
	}
	// Note that this doesn't have an ID/USERNAME to set password. This table only ever
	// has 1 row (maybe that will change in the future? if this this will not work)
	err = it.PasswordStore.Set(string(hash))
	if err != nil {
		return err
	}

	r := &SuccessResponse{Status: "ok", Data: nil}
	return c.JSON(http.StatusOK, r)
}

type ValidateUserRequestBody struct {
	Token core.Token `json:"token"`
}

type ValidateUserData struct {
	Validated bool `json:"validated"`
}
type ValidateUserResponse struct {
	SuccessResponse
	Data ValidateUserData `json:"data"`
}

func (it *RouteHandler) ValidateUser(c echo.Context) error {
	req := new(ValidateUserRequestBody)
	if err := c.Bind(req); err != nil {
		return err
	}
	val := it.authenticateUser(c, &req.Token)
	if !val {
		r := &FailureResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	r := &ValidateUserResponse{SuccessResponse: SuccessResponse{Status: "ok"}, Data: ValidateUserData{Validated: true}}
	return c.JSON(http.StatusOK, r)
}

func (it *RouteHandler) authenticateUser(c echo.Context, token *core.Token) bool {
	if *token == "" {
		*token = c.Request().Header.Get("x-actual-token")
	}
	res, err := it.TokenStore.Has(*token)
	if err != nil {
		return false
	}
	return res
}
