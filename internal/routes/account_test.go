package routes_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/routes"
	"github.com/nathanjisaac/actual-server-go/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupAccountTestHandler(body string, pstore core.PasswordStore, tstore core.TokenStore) (*routes.RouteHandler, echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &routes.RouteHandler{
		Config: core.Config{
			Mode: core.Development,
		},
		FileStore:     nil,
		PasswordStore: nil,
		TokenStore:    nil,
	}
	if pstore != nil && tstore == nil {
		h = &routes.RouteHandler{
			Config: core.Config{
				Mode: core.Development,
			},
			FileStore:     nil,
			PasswordStore: pstore,
			TokenStore:    nil,
		}
	} else if pstore == nil && tstore != nil {
		h = &routes.RouteHandler{
			Config: core.Config{
				Mode: core.Development,
			},
			FileStore:     nil,
			PasswordStore: nil,
			TokenStore:    tstore,
		}
	} else if pstore != nil && tstore != nil {
		h = &routes.RouteHandler{
			Config: core.Config{
				Mode: core.Development,
			},
			FileStore:     nil,
			PasswordStore: pstore,
			TokenStore:    tstore,
		}
	}
	return h, c, rec
}

func TestNeedsBootstrap(t *testing.T) {
	t.Run("given no passwords then return not bootstrapped", func(t *testing.T) {
		store := memory.NewPasswordStore()
		h, c, rec := setupAccountTestHandler("", store, nil)

		var res routes.NeedsBootstrapResponse
		err := h.NeedsBootstrap(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

		assert.Equal(t, false, res.Data.Bootstrapped)
	})

	t.Run("given a password then return bootstrapped", func(t *testing.T) {
		store := memory.NewPasswordStore()
		h, c, rec := setupAccountTestHandler("", store, nil)

		err := store.Add("password")
		assert.NoError(t, err)

		var res routes.NeedsBootstrapResponse
		err = h.NeedsBootstrap(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

		assert.Equal(t, true, res.Data.Bootstrapped)
	})
}

func TestBootstrap(t *testing.T) {
	t.Run("given empty body password then returns error", func(t *testing.T) {
		h, c, rec := setupAccountTestHandler("", nil, nil)

		var res routes.ErrorResponse
		err := h.Bootstrap(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "invalid-password", res.Reason)
	})

	t.Run("given already bootstrapped then returns error", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		h, c, rec := setupAccountTestHandler(`{"password":"pass"}`, pStore, nil)

		err := pStore.Add("password")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.Bootstrap(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "already-bootstrapped", res.Reason)
	})

	t.Run("given not bootstrapped then returns token", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		h, c, rec := setupAccountTestHandler(`{"password":"pass"}`, pStore, tStore)

		var res routes.BootstrapResponse
		err := h.Bootstrap(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		_, err = uuid.Parse(res.Data.Token)
		assert.NoError(t, err)
	})
}

func TestLogin(t *testing.T) {
	t.Run("given empty password then returns no token", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		h, c, rec := setupAccountTestHandler("", pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(uuid.NewString())
		assert.NoError(t, err)

		var res routes.LoginFailResponse
		err = h.Login(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.NotEqual(t, nil, res.Data)
		assert.Equal(t, nil, res.Data.Token)
	})

	t.Run("given not bootstrapped then returns no token", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		h, c, rec := setupAccountTestHandler("", pStore, tStore)

		var res routes.LoginFailResponse
		err := h.Login(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.NotEqual(t, nil, res.Data)
		assert.Equal(t, nil, res.Data.Token)
	})

	t.Run("given wrong password then returns no token", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		h, c, rec := setupAccountTestHandler(`{"password":"pass"}`, pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(uuid.NewString())
		assert.NoError(t, err)

		var res routes.LoginFailResponse
		err = h.Login(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.NotEqual(t, nil, res.Data)
		assert.Equal(t, nil, res.Data.Token)
	})

	t.Run("given correct password then returns token", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		h, c, rec := setupAccountTestHandler(`{"password":"password123"}`, pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		token := uuid.NewString()
		err = tStore.Add(token)
		assert.NoError(t, err)

		var res routes.LoginSuccessResponse
		err = h.Login(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.NotEqual(t, nil, res.Data)
		assert.Equal(t, token, res.Data.Token)
	})
}

func TestChangePassword(t *testing.T) {
	t.Run("given no token in body/header then returns error", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler("", pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.ChangePassword(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given empty password with token in body then returns error", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler(fmt.Sprintf(`{"token":"%s"}`, token), pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.ChangePassword(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "invalid-password", res.Reason)
	})

	t.Run("given empty password with token in header then returns error", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler("", pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", token)

		var res routes.ErrorResponse
		err = h.ChangePassword(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "invalid-password", res.Reason)
	})

	t.Run("given password with token in body then returns success", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler(fmt.Sprintf(`{"token":"%s","password":"password456"}`, token), pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)

		var res routes.SuccessResponse
		err = h.ChangePassword(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.Equal(t, nil, res.Data)
	})

	t.Run("given password with token in header then returns success", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler(`{"password":"password456"}`, pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", token)

		var res routes.SuccessResponse
		err = h.ChangePassword(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.Equal(t, nil, res.Data)
	})
}

func TestValidateUser(t *testing.T) {
	t.Run("given no token in body/header then returns error", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler("", pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.ValidateUser(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given invalid token in body then returns error", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler(fmt.Sprintf(`{"token":"%s"}`, token), pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(uuid.NewString())
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.ValidateUser(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given invalid token in header then returns error", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler("", pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", uuid.NewString())

		var res routes.ErrorResponse
		err = h.ValidateUser(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given valid token in body then returns success", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler(fmt.Sprintf(`{"token":"%s"}`, token), pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)

		var res routes.ValidateUserResponse
		err = h.ValidateUser(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.NotEqual(t, nil, res.Data)
		assert.Equal(t, true, res.Data.Validated)
	})

	t.Run("given valid token in header then returns success", func(t *testing.T) {
		pStore := memory.NewPasswordStore()
		tStore := memory.NewTokenStore()
		token := uuid.NewString()
		h, c, rec := setupAccountTestHandler("", pStore, tStore)

		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		assert.NoError(t, err)
		err = pStore.Add(string(hash))
		assert.NoError(t, err)
		err = tStore.Add(token)
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", token)

		var res routes.ValidateUserResponse
		err = h.ValidateUser(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.NotEqual(t, nil, res.Data)
		assert.Equal(t, true, res.Data.Validated)
	})
}
