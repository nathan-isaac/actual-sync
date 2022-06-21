package routes

import (
	"github.com/nathanjisaac/actual-server-go/internal/core"
)

type RouteHandler struct {
	Config        core.Config
	FileStore     core.FileStore
	PasswordStore core.PasswordStore
	TokenStore    core.TokenStore
}

type ErrorResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
