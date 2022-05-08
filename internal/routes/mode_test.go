package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetMode(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &RouteHandler{Config: core.Config{Mode: core.Development}}

	if assert.NoError(t, h.GetMode(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "development", rec.Body.String())
	}
}
