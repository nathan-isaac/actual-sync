package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/routes"
	"github.com/stretchr/testify/assert"
)

func TestGetMode(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()

	h := &routes.RouteHandler{Config: core.Config{Mode: core.Development}}

	if c := e.NewContext(req, rec); assert.NoError(t, h.GetMode(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "development", rec.Body.String())
	}
}
