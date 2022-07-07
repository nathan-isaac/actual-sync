//nolint: dupl // Disabling dupl for tests. It detects similar testcases for different tests.
package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/routes"
	"github.com/nathanjisaac/actual-server-go/internal/storage/sqlite"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func setupSyncTestHandler(body string, tstore core.TokenStore, fstore core.FileStore) (
	*routes.RouteHandler,
	echo.Context,
	*httptest.ResponseRecorder,
) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	config := core.Config{
		Mode: core.Development,
	}
	h := &routes.RouteHandler{
		Config:        config,
		FileStore:     nil,
		PasswordStore: nil,
		TokenStore:    nil,
	}
	switch fstore {
	case nil:
		if tstore != nil {
			h = &routes.RouteHandler{
				Config:        config,
				FileStore:     nil,
				PasswordStore: nil,
				TokenStore:    tstore,
			}
		}
	default:
		switch tstore {
		case nil:
			h = &routes.RouteHandler{
				Config:        config,
				FileStore:     fstore,
				PasswordStore: nil,
				TokenStore:    nil,
			}
		default:
			h = &routes.RouteHandler{
				Config:        config,
				FileStore:     fstore,
				PasswordStore: nil,
				TokenStore:    tstore,
			}
		}
	}
	return h, c, rec
}

func setupSyncTestFileHandler(body []byte, tstore core.TokenStore, fstore core.FileStore, fileID string) (
	*routes.RouteHandler,
	echo.Context,
	*httptest.ResponseRecorder,
) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment;filename=%s", fileID))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fs := afero.NewMemMapFs()
	config := core.Config{
		Mode:       core.Development,
		FileSystem: fs,
		UserFiles:  "",
	}
	h := &routes.RouteHandler{
		Config:        config,
		FileStore:     nil,
		PasswordStore: nil,
		TokenStore:    nil,
	}
	switch fstore {
	case nil:
		if tstore != nil {
			h = &routes.RouteHandler{
				Config:        config,
				FileStore:     nil,
				PasswordStore: nil,
				TokenStore:    tstore,
			}
		}
	default:
		switch tstore {
		case nil:
			h = &routes.RouteHandler{
				Config:        config,
				FileStore:     fstore,
				PasswordStore: nil,
				TokenStore:    nil,
			}
		default:
			h = &routes.RouteHandler{
				Config:        config,
				FileStore:     fstore,
				PasswordStore: nil,
				TokenStore:    tstore,
			}
		}
	}
	return h, c, rec
}

func TestUserCreateKey(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"fileId":"1","keyId":"2","keySalt":"3","testContent":"4"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.UserCreateKey(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given token and no valid file returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, _ := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		err = h.UserCreateKey(c)
		assert.Error(t, err)
		assert.Equal(t, "no record updated", err.Error())
	})

	t.Run("given token and valid file returns success", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(
			`{"token":"token123","fileId":"f1","keyId":"2","keySalt":"3","testContent":"4"}`,
			tstore,
			fstore,
		)

		err = tstore.Add("token123")
		assert.NoError(t, err)
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)

		var res routes.SuccessResponse
		err = h.UserCreateKey(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)

		file, err := fstore.ForID("f1")
		assert.NoError(t, err)
		assert.Equal(t, "2", file.EncryptKeyID)
		assert.Equal(t, "3", file.EncryptSalt)
		assert.Equal(t, "4", file.EncryptTest)
	})
}

func TestUserGetKey(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"fileId":"1"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.UserGetKey(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given token and no valid file returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		err = h.UserGetKey(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "file-not-found", rec.Body.String())
	})

	t.Run("given token and valid file returns success", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(
			`{"token":"token123","fileId":"f1","keyId":"2","keySalt":"3","testContent":"4"}`,
			tstore,
			fstore,
		)

		err = tstore.Add("token123")
		assert.NoError(t, err)
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		var res routes.UserGetKeyResponse
		err = h.UserGetKey(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.Equal(t, "keyid", res.Data.EncryptKeyID)
		assert.Equal(t, "salt", res.Data.EncryptSalt)
		assert.Equal(t, "test", res.Data.EncryptTest)
	})
}

func TestResetUserFile(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"fileId":"1"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.ResetUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given token and no valid file returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		err = h.ResetUserFile(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "User or file not found", rec.Body.String())
	})

	t.Run("given token and valid file returns success", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123","fileId":"f1"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		var res routes.SuccessResponse
		err = h.ResetUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)

		file, err := fstore.ForID("f1")
		assert.NoError(t, err)
		assert.Equal(t, "", file.GroupID)
	})
}

func TestUpdateUserFileName(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"fileId":"1","name":"budgetnew"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.UpdateUserFileName(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given token and no valid file returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		err = h.UpdateUserFileName(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "User or file not found", rec.Body.String())
	})

	t.Run("given token and valid file returns success", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123","fileId":"f1","name":"budgetnew"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		var res routes.SuccessResponse
		err = h.UpdateUserFileName(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)

		file, err := fstore.ForID("f1")
		assert.NoError(t, err)
		assert.Equal(t, "budgetnew", file.Name)
	})
}

func TestUserFileInfo(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"fileId":"1"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.UserFileInfo(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given token and no files returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-file-id", "f1")

		var res routes.ErrorResponse
		err = h.UserFileInfo(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "User or file not found", res.Reason)
	})

	t.Run("given token and file exists returns info without encrypt meta", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-file-id", "f1")

		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		var res routes.UserFileInfoResponse
		err = h.UserFileInfo(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)

		assert.Equal(t, "f1", res.Data.FileID)
		assert.Equal(t, "g1", res.Data.GroupID)
		assert.Equal(t, "budget", res.Data.Name)
		assert.Equal(t, false, res.Data.Deleted)
	})

	t.Run("given token and file exists returns info with encrypt meta", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-file-id", "f1")

		err = fstore.Add(&core.NewFile{
			FileID:      "f1",
			GroupID:     "g1",
			SyncVersion: 2,
			EncryptMeta: `{"keyId":"keyidMeta"}`,
			Name:        "budget",
		})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		var res routes.UserFileInfoWithMetaResponse
		err = h.UserFileInfo(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)

		assert.Equal(t, "f1", res.Data.FileID)
		assert.Equal(t, "g1", res.Data.GroupID)
		assert.Equal(t, "keyidMeta", res.Data.EncryptMeta.KeyID)
		assert.Equal(t, "budget", res.Data.Name)
		assert.Equal(t, false, res.Data.Deleted)
	})
}

func TestListUserFiles(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"fileId":"1"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.ListUserFiles(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given token and no files returns empty array", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ListFilesResponse
		err = h.ListUserFiles(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.Equal(t, 0, len(res.Data))
	})

	t.Run("given token and two valid file returns array", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123","fileId":"f1"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		err = fstore.Add(&core.NewFile{FileID: "f2", GroupID: "g2", SyncVersion: 2, EncryptMeta: "abc2", Name: "budget2"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f2", "salt2", "keyid2", "test2")
		assert.NoError(t, err)

		var res routes.ListFilesResponse
		err = h.ListUserFiles(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.Equal(t, 2, len(res.Data))

		assert.Equal(t, "f1", res.Data[0].FileID)
		assert.Equal(t, "g1", res.Data[0].GroupID)
		assert.Equal(t, "keyid", res.Data[0].EncryptKeyID)
		assert.Equal(t, "budget", res.Data[0].Name)
		assert.Equal(t, false, res.Data[0].Deleted)

		assert.Equal(t, "f2", res.Data[1].FileID)
		assert.Equal(t, "g2", res.Data[1].GroupID)
		assert.Equal(t, "keyid2", res.Data[1].EncryptKeyID)
		assert.Equal(t, "budget2", res.Data[1].Name)
		assert.Equal(t, false, res.Data[1].Deleted)
	})
}

func TestUploadUserFIle(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestFileHandler([]byte{}, tstore, fstore, "")

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.UploadUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given logged in and no files then adds file and returns success", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestFileHandler([]byte("testing"), tstore, fstore, "f1")

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", "token123")
		c.Request().Header.Set("x-actual-name", "budget")
		c.Request().Header.Set("x-actual-file-id", "f1")
		c.Request().Header.Set("x-actual-group-id", "g1")
		c.Request().Header.Set("x-actual-encrypt-meta", `{"keyId": "keyid"}`)
		c.Request().Header.Set("x-actual-format", "2")

		var res routes.UploadUserFileResponse
		err = h.UploadUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		_, err = uuid.Parse(res.GroupID)
		assert.NoError(t, err)

		count, err := fstore.Count()
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		fs := h.Config.FileSystem
		result, err := afero.FileContainsBytes(fs, filepath.Join(h.Config.UserFiles, "f1.blob"), []byte("testing"))
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("given logged in, old groupid and file exists then returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestFileHandler([]byte("testing"), tstore, fstore, "f1")

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", "token123")
		c.Request().Header.Set("x-actual-name", "budget")
		c.Request().Header.Set("x-actual-file-id", "f1")
		c.Request().Header.Set("x-actual-group-id", "g2")
		c.Request().Header.Set("x-actual-encrypt-meta", `{"keyId": "keyid"}`)
		c.Request().Header.Set("x-actual-format", "2")
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		err = h.UploadUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "file-has-reset", rec.Body.String())
	})

	t.Run("given logged in, old encypt keyid and file exists then returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestFileHandler([]byte("testing"), tstore, fstore, "f1")

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", "token123")
		c.Request().Header.Set("x-actual-name", "budget")
		c.Request().Header.Set("x-actual-file-id", "f1")
		c.Request().Header.Set("x-actual-group-id", "g1")
		c.Request().Header.Set("x-actual-encrypt-meta", `{"keyId": "keyid"}`)
		c.Request().Header.Set("x-actual-format", "2")
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid2", "test")
		assert.NoError(t, err)

		err = h.UploadUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "file-has-new-key", rec.Body.String())
	})

	t.Run("given logged in and file exists then returns success", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestFileHandler([]byte("testing"), tstore, fstore, "f1")

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", "token123")
		c.Request().Header.Set("x-actual-name", "budgetnew")
		c.Request().Header.Set("x-actual-file-id", "f1")
		c.Request().Header.Set("x-actual-group-id", "g1")
		c.Request().Header.Set("x-actual-encrypt-meta", `{"keyId": "keyid"}`)
		c.Request().Header.Set("x-actual-format", "3")
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		var res routes.UploadUserFileResponse
		err = h.UploadUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)
		assert.Equal(t, "g1", res.GroupID)

		file, err := fstore.ForID("f1")
		assert.NoError(t, err)
		assert.Equal(t, int16(3), file.SyncVersion)
		assert.Equal(t, "budgetnew", file.Name)

		fs := h.Config.FileSystem
		result, err := afero.FileContainsBytes(fs, filepath.Join(h.Config.UserFiles, "f1.blob"), []byte("testing"))
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})
}

func TestDownloadUserFIle(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestFileHandler([]byte{}, tstore, fstore, "")

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.DownloadUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given logged in and no files then returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestFileHandler([]byte("testing"), tstore, fstore, "f1")

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", "token123")
		c.Request().Header.Set("x-actual-file-id", "f1")

		err = h.DownloadUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "User or file not found", rec.Body.String())
	})

	t.Run("given logged in and file exists then returns success", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestFileHandler([]byte("testing"), tstore, fstore, "f1")

		err = tstore.Add("token123")
		assert.NoError(t, err)
		c.Request().Header.Set("x-actual-token", "token123")
		c.Request().Header.Set("x-actual-file-id", "f1")
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		out, err := h.Config.FileSystem.Create(filepath.Join(h.Config.UserFiles, "f1.blob"))
		assert.NoError(t, err)
		defer out.Close()
		_, err = io.Copy(out, bytes.NewReader([]byte("testing")))
		assert.NoError(t, err)

		err = h.DownloadUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, []byte("testing"), rec.Body.Bytes())
	})
}

func TestDeleteUserFile(t *testing.T) {
	t.Run("given no token in returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"fileId":"1"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		var res routes.ErrorResponse
		err = h.DeleteUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "error", res.Status)
		assert.Equal(t, "auth-error", res.Reason)
	})

	t.Run("given token and no valid file returns error", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)

		err = h.DeleteUserFile(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "User or file not found", rec.Body.String())
	})

	t.Run("given token and valid file returns success", func(t *testing.T) {
		db, err := sqlite.NewAccountConnection(":memory:")
		assert.NoError(t, err)
		defer db.Close()
		tstore := sqlite.NewTokenStore(db)
		fstore := sqlite.NewFileStore(db)
		h, c, rec := setupSyncTestHandler(`{"token":"token123","fileId":"f1"}`, tstore, fstore)

		err = tstore.Add("token123")
		assert.NoError(t, err)
		err = fstore.Add(&core.NewFile{FileID: "f1", GroupID: "g1", SyncVersion: 2, EncryptMeta: "abc", Name: "budget"})
		assert.NoError(t, err)
		err = fstore.UpdateEncryption("f1", "salt", "keyid", "test")
		assert.NoError(t, err)

		var res routes.SuccessResponse
		err = h.DeleteUserFile(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(t, "ok", res.Status)

		file, err := fstore.ForID("f1")
		assert.NoError(t, err)
		assert.Equal(t, true, file.Deleted)
	})
}
