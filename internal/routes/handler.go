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
