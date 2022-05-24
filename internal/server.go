package internal

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/routes"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
)

func StartServer(config core.Config, BuildDirectory embed.FS, headless bool) {
	e := echo.New()
	e.HideBanner = true

	if !headless {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "node_modules/@actual-app/web/build",
			HTML5:      true,
			Filesystem: http.FS(BuildDirectory),
		}))
	}

	conn, err := sqlite.NewAccountConnection(filepath.Join(config.ServerFiles, "account.sqlite"))
	if err != nil {
		e.Logger.Fatal(err)
	}

	handler := routes.RouteHandler{
		Config:        config,
		FileStore:     sqlite.NewFileStore(conn),
		TokenStore:    sqlite.NewTokenStore(conn),
		PasswordStore: sqlite.NewPasswordStore(conn),
	}
	e.GET("/mode", handler.GetMode)

	account := e.Group("/account")
	account.GET("/needs-bootstrap", handler.NeedsBootstrap)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", config.Hostname, config.Port)))
}
