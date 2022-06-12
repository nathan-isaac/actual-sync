package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage"
)

type UserCreateKeyRequestBody struct {
	FileId      string `json:"fileId"`
	KeyId       string `json:"keyId"`
	KeySalt     string `json:"keySalt"`
	TestContent string `json:"testContent"`
	Token       string `json:"token"`
}

func (it *RouteHandler) UserCreateKey(c echo.Context) error {
	req := new(UserCreateKeyRequestBody)
	if err := c.Bind(req); err != nil {
		c.Echo().Logger.Error(err)
		return err
	}

	val := it.authenticateUser(c, req.Token)
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	err := it.FileStore.UpdateEncryption(req.FileId, req.KeySalt, req.KeyId, req.TestContent)
	if err != nil {
		c.Echo().Logger.Error(err)
		return err
	}

	r := &SuccessResponse{Status: "ok"}
	return c.JSON(http.StatusOK, r)
}

type UserGetKeyRequestBody struct {
	FileId string `json:"fileId"`
	Token  string `json:"token"`
}

type UserGetKeyResponse struct {
	SuccessResponse
	Data UserGetKeyResponseData `json:"data"`
}

type UserGetKeyResponseData struct {
	EncryptKeyId string `json:"id"`
	EncryptSalt  string `json:"salt"`
	EncryptTest  string `json:"test"`
}

func (it *RouteHandler) UserGetKey(c echo.Context) error {
	req := new(UserGetKeyRequestBody)
	if err := c.Bind(req); err != nil {
		c.Echo().Logger.Error(err)
		return err
	}

	val := it.authenticateUser(c, req.Token)
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	file, err := it.FileStore.ForId(req.FileId)
	if err != nil {
		if err == storage.ErrorRecordNotFound {
			return c.String(http.StatusBadRequest, "file-not-found")
		}
		c.Echo().Logger.Error(err)
		return err
	}

	r := &UserGetKeyResponse{
		SuccessResponse: SuccessResponse{Status: "ok"},
		Data: UserGetKeyResponseData{
			EncryptKeyId: file.EncryptKeyId,
			EncryptSalt:  file.EncryptSalt,
			EncryptTest:  file.EncryptTest,
		},
	}
	return c.JSON(http.StatusOK, r)
}

func (it *RouteHandler) ResetUserFile(c echo.Context) error {
	req := new(UserGetKeyRequestBody)
	if err := c.Bind(req); err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	val := it.authenticateUser(c, req.Token)
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	err := it.FileStore.ClearGroup(req.FileId)
	if err != nil {
		if err == storage.ErrorNoRecordUpdated {
			return c.String(http.StatusBadRequest, "User or file not found")
		}
		c.Echo().Logger.Error(err)
		return err
	}

	r := &SuccessResponse{Status: "ok"}
	return c.JSON(http.StatusOK, r)
}

type UpdateUserFileNameRequestBody struct {
	FileId string `json:"fileId"`
	Name   string `json:"name"`
	Token  string `json:"token"`
}

func (it *RouteHandler) UpdateUserFileName(c echo.Context) error {
	req := new(UpdateUserFileNameRequestBody)
	if err := c.Bind(req); err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	val := it.authenticateUser(c, req.Token)
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	err := it.FileStore.UpdateName(req.FileId, req.Name)
	if err != nil {
		if err == storage.ErrorNoRecordUpdated {
			return c.String(http.StatusBadRequest, "User or file not found")
		}
		c.Echo().Logger.Error(err)
		return err
	}

	r := &SuccessResponse{Status: "ok"}
	return c.JSON(http.StatusOK, r)
}

type encryptMetaType struct {
	KeyId string `json:"keyId"`
}

type UserFileInfoWithMetaResponse struct {
	SuccessResponse
	Data FileInfoWithMetaResponseData `json:"data"`
}

type FileInfoWithMetaResponseData struct {
	Name        string          `json:"name"`
	FileId      string          `json:"fileId"`
	GroupId     string          `json:"groupId"`
	EncryptMeta encryptMetaType `json:"encryptMeta"`
	Deleted     bool            `json:"deleted"`
}

type UserFileInfoResponse struct {
	SuccessResponse
	Data FileInfoResponseData `json:"data"`
}

type FileInfoResponseData struct {
	Name    string `json:"name"`
	FileId  string `json:"fileId"`
	GroupId string `json:"groupId"`
	Deleted bool   `json:"deleted"`
}

func (it *RouteHandler) UserFileInfo(c echo.Context) error {
	req := new(UserGetKeyRequestBody)
	req.FileId = c.Request().Header.Get("x-actual-file-id")
	if err := c.Bind(req); err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	val := it.authenticateUser(c, req.Token)
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	file, err := it.FileStore.ForIdAndDelete(req.FileId, false)
	if err != nil {
		if err == storage.ErrorRecordNotFound {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Status: "error", Reason: "User or file not found"})
		}
		c.Echo().Logger.Error(err)
		return err
	}

	if file.EncryptMeta != "" {
		var meta encryptMetaType
		err = json.Unmarshal([]byte(file.EncryptMeta), &meta)
		if err != nil {
			c.Echo().Logger.Error(err)
			return err
		}
		r := UserFileInfoWithMetaResponse{
			SuccessResponse: SuccessResponse{Status: "ok"},
			Data:            FileInfoWithMetaResponseData{Name: file.Name, FileId: file.FileId, GroupId: file.GroupId, EncryptMeta: meta, Deleted: file.Deleted},
		}
		return c.JSON(http.StatusOK, r)
	}

	r := UserFileInfoResponse{
		SuccessResponse: SuccessResponse{Status: "ok"},
		Data:            FileInfoResponseData{Name: file.Name, FileId: file.FileId, GroupId: file.GroupId, Deleted: file.Deleted},
	}
	return c.JSON(http.StatusOK, r)
}

type TokenRequestBody struct {
	Token string `json:"token"`
}

type ListFilesResponse struct {
	SuccessResponse
	Data []FileResponseData `json:"data"`
}

type FileResponseData struct {
	Name         string `json:"name"`
	FileId       string `json:"fileId"`
	GroupId      string `json:"groupId"`
	EncryptKeyId string `json:"encryptKeyIid"`
	Deleted      bool   `json:"deleted"`
}

func (it *RouteHandler) ListUserFiles(c echo.Context) error {
	req := new(TokenRequestBody)
	if err := c.Bind(req); err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	val := it.authenticateUser(c, req.Token)
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	files, err := it.FileStore.All()
	if err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	filesRes := make([]FileResponseData, 0)
	for _, file := range files {
		filesRes = append(filesRes, FileResponseData{Name: file.Name, FileId: file.FileId, GroupId: file.GroupId, EncryptKeyId: file.EncryptKeyId, Deleted: file.Deleted})
	}

	r := &ListFilesResponse{
		SuccessResponse: SuccessResponse{Status: "ok"},
		Data:            filesRes,
	}
	return c.JSON(http.StatusOK, r)
}

type UploadUserFileResponse struct {
	SuccessResponse
	GroupId string `json:"groupId"`
}

func (it *RouteHandler) UploadUserFile(c echo.Context) error {
	val := it.authenticateUser(c, "")
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	name, err := url.PathUnescape(c.Request().Header.Get("x-actual-name"))
	if err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	fileId := c.Request().Header.Get("x-actual-file-id")
	groupId := c.Request().Header.Get("x-actual-group-id")
	encryptMeta := c.Request().Header.Get("x-actual-encrypt-meta")
	syncFormatVersion, err := strconv.ParseInt(c.Request().Header.Get("x-actual-format"), 10, 16)
	if err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	keyId := ""
	if encryptMeta != "" {
		var jsonData encryptMetaType
		err := json.Unmarshal([]byte(encryptMeta), &jsonData)
		if err != nil {
			c.Echo().Logger.Error(err)
			return err
		}
		keyId = jsonData.KeyId
	}

	file, err := it.FileStore.ForId(fileId)
	fileExists := false
	if err != nil && err != storage.ErrorRecordNotFound {
		c.Echo().Logger.Error(err)
		return err
	}
	if err != storage.ErrorRecordNotFound {
		fileExists = true
		// The uploading file is part of an old group, so reject
		// it. All of its internal sync state is invalid because its
		// old. The sync state has been reset, so user needs to
		// either reset again or download from the current group.
		if groupId != file.GroupId {
			return c.String(http.StatusBadRequest, "file-has-reset")
		}

		// The key that the file is encrypted with is different than
		// the current registered key. All data must always be
		// encrypted with the registered key for consistency. Key
		// changes always necessitate a sync reset, which means this
		// upload is trying to overwrite another reset. That might
		// be be fine, but since we definitely cannot accept a file
		// encrypted with the wrong key, we bail and suggest the
		// user download the latest file.
		if keyId != file.EncryptKeyId {
			return c.String(http.StatusBadRequest, "file-has-new-key")
		}
	}

	out, err := it.Config.FileSystem.Create(filepath.Join(it.Config.UserFiles, fmt.Sprintf("%s.blob", fileId)))
	if err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, c.Request().Body)
	if err != nil {
		c.Echo().Logger.Error(err)
		return err
	}

	if !fileExists {
		// Its new
		uuid, err := uuid.NewRandom()
		if err != nil {
			c.Echo().Logger.Error(err)
			return err
		}
		groupId = uuid.String()

		err = it.FileStore.Add(&core.NewFile{FileId: fileId, GroupId: groupId, SyncVersion: int16(syncFormatVersion), EncryptMeta: encryptMeta, Name: name})
		if err != nil {
			c.Echo().Logger.Error(err)
			return err
		}

		r := UploadUserFileResponse{
			SuccessResponse: SuccessResponse{Status: "ok"},
			GroupId:         groupId,
		}
		return c.JSON(http.StatusOK, r)
	} else {
		if groupId == "" {
			// Sync state was reset. Create new group
			uuid, err := uuid.NewRandom()
			if err != nil {
				c.Echo().Logger.Error(err)
				return err
			}
			groupId = uuid.String()

			err = it.FileStore.UpdateGroup(fileId, groupId)
			if err != nil {
				c.Echo().Logger.Error(err)
				return err
			}
		}

		// Regardless, update properties
		err = it.FileStore.Update(fileId, int16(syncFormatVersion), encryptMeta, name)
		if err != nil {
			c.Echo().Logger.Error(err)
			return err
		}

		r := UploadUserFileResponse{
			SuccessResponse: SuccessResponse{Status: "ok"},
			GroupId:         groupId,
		}
		return c.JSON(http.StatusOK, r)
	}
}

func (it *RouteHandler) DownloadUserFile(c echo.Context) error {
	val := it.authenticateUser(c, "")
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	fileId := c.Request().Header.Get("x-actual-file-id")

	_, err := it.FileStore.ForIdAndDelete(fileId, false)
	if err != nil {
		if err == storage.ErrorRecordNotFound {
			return c.String(http.StatusBadRequest, "User or file not found")
		}
		c.Echo().Logger.Error(err)
		return err
	}

	fs := it.Config.FileSystem
	file, err := fs.Open(filepath.Join(it.Config.UserFiles, fmt.Sprintf("%s.blob", fileId)))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading files")
	}
	defer file.Close()
	finfo, err := file.Stat()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading files")
	}
	fileBlob := make([]byte, finfo.Size())
	_, err = file.Read(fileBlob)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading files")
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", fileId))
	return c.Blob(http.StatusOK, echo.MIMEOctetStream, fileBlob)
}

func (it *RouteHandler) DeleteUserFile(c echo.Context) error {
	req := new(UserGetKeyRequestBody)
	if err := c.Bind(req); err != nil {
		c.Echo().Logger.Error(err)
		return err
	}
	val := it.authenticateUser(c, req.Token)
	if !val {
		r := &ErrorResponse{
			Status: "error",
			Reason: "auth-error",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	err := it.FileStore.Delete(req.FileId)
	if err != nil {
		if err == storage.ErrorNoRecordUpdated {
			return c.String(http.StatusBadRequest, "User or file not found")
		}
		c.Echo().Logger.Error(err)
		return err
	}

	r := &SuccessResponse{Status: "ok"}
	return c.JSON(http.StatusOK, r)
}