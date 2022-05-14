package routes

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNeedsBootstrap(t *testing.T) {
	t.Run("given no passwords then return not bootstrapped", func(t *testing.T) {
		store := memory.New()

		res := newResponse(t, store)

		assert.Equal(t, false, res.Data.Bootstrapped)
	})

	t.Run("given a password then return bootstrapped", func(t *testing.T) {
		store := memory.New()
		err := store.Add("password")
		assert.NoError(t, err)

		res := newResponse(t, store)

		assert.Equal(t, true, res.Data.Bootstrapped)
	})
}

func newResponse(t *testing.T, store core.PasswordStore) NeedsBootstrapResponse {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &RouteHandler{
		Config: core.Config{
			Mode: core.Development,
		},
		FileStore:     nil,
		PasswordStore: store,
		TokenStore:    nil,
	}

	var response NeedsBootstrapResponse

	err := h.NeedsBootstrap(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	return response
}
