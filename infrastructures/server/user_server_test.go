package server_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/server"
	"github.com/wisle25/be-template/infrastructures/services"
	"github.com/wisle25/be-template/tests/db_helper"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUserHTTP(t *testing.T) {
	// Prepare
	config := commons.LoadConfig("../..")
	db := services.ConnectDB(config)
	userHelperDB := &db_helper.UserHelperDB{
		DB: db,
	}
	defer userHelperDB.CleanUserDB()

	var accessTokenCookie string
	var refreshTokenCookie string

	app := server.CreateServer(config)

	t.Run("Register User", func(t *testing.T) {
		// Arrange
		payload := map[string]string{
			"username":        "user",
			"email":           "user@example.com",
			"password":        "password",
			"confirmPassword": "password",
		}
		body, _ := json.Marshal(payload)

		// Action
		req := httptest.NewRequest(
			"POST",
			"/users",
			bytes.NewBuffer(body),
		)
		req.Header.Set("Content-Type", "application/json")

		res, _ := app.Test(req)

		// Assert
		resBody, _ := io.ReadAll(res.Body)
		var resMap map[string]string
		_ = json.Unmarshal(resBody, &resMap)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.NotNil(t, resMap["id"])
	})

	t.Run("Login User", func(t *testing.T) {
		// Arrange
		payload := map[string]string{
			"identity": "user",
			"password": "password",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(
			"POST",
			"/auths",
			bytes.NewBuffer(body),
		)
		req.Header.Set("Content-Type", "application/json")

		// Action
		res, _ := app.Test(req)

		// Save cookie for refresh token test
		for _, cookie := range res.Cookies() {
			if cookie.Name == "refresh_token" {
				refreshTokenCookie = cookie.Value
			} else if cookie.Name == "access_token" {
				accessTokenCookie = cookie.Value
			}
		}

		// Assert
		resBody, _ := io.ReadAll(res.Body)
		var resMap map[string]string
		_ = json.Unmarshal(resBody, &resMap)

		assert.Equal(t, 3, len(res.Cookies()))
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "success", resMap["status"])
		assert.NotEmpty(t, refreshTokenCookie)
	})

	t.Run("Refresh Token", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(
			"PUT",
			"/auths",
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: refreshTokenCookie,
		})

		// Action
		res, _ := app.Test(req)

		// Assert
		resBody, _ := io.ReadAll(res.Body)
		var resMap map[string]string
		_ = json.Unmarshal(resBody, &resMap)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "success", resMap["status"])
	})

	t.Run("Logout User", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(
			http.MethodDelete,
			"/auths",
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: accessTokenCookie,
		})
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: refreshTokenCookie,
		})

		// Action
		res, _ := app.Test(req)

		// Assert
		resBody, _ := io.ReadAll(res.Body)
		var resMap map[string]string
		_ = json.Unmarshal(resBody, &resMap)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "success", resMap["status"])
	})

	t.Run("Get User By Id", func(t *testing.T) {
		// Arrange
		expectedUser := &userHelperDB.GetUsers()[0]

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/"+expectedUser.Id,
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: accessTokenCookie,
		})
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: refreshTokenCookie,
		})

		// Action
		res, _ := app.Test(req)

		// Assert
		resBody, _ := io.ReadAll(res.Body)
		var resMap map[string]interface{}
		_ = json.Unmarshal(resBody, &resMap)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "success", resMap["status"])

		// Access the data part
		dataMap := resMap["data"].(map[string]interface{})

		assert.Equal(t, expectedUser.Id, dataMap["Id"].(string))
		assert.Equal(t, expectedUser.Username, dataMap["Username"].(string))
		assert.Equal(t, expectedUser.Email, dataMap["Email"].(string))
		assert.Equal(t, expectedUser.AvatarLink, dataMap["AvatarLink"].(string))
	})

	t.Run("Update User By Id", func(t *testing.T) {
		// Arrange
		expectedUser := &userHelperDB.GetUsers()[0]
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		_ = writer.WriteField("username", "updateduser")
		_ = writer.WriteField("email", "updateduser@example.com")
		_ = writer.WriteField("password", "newpassword")
		_ = writer.WriteField("confirmPassword", "newpassword")

		file, err := os.Open(filepath.Join("tests", "avatar.png"))
		assert.Nil(t, err)
		part, err := writer.CreateFormFile("avatar", "avatar.png")
		assert.Nil(t, err)
		_, err = io.Copy(part, file)
		assert.Nil(t, err)
		file.Close()

		writer.Close()

		req := httptest.NewRequest(
			http.MethodPut,
			"/users/"+expectedUser.Id,
			&buf,
		)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: accessTokenCookie,
		})
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: refreshTokenCookie,
		})

		// Action
		res, _ := app.Test(req)

		// Assert
		resBody, _ := io.ReadAll(res.Body)
		var resMap map[string]string
		_ = json.Unmarshal(resBody, &resMap)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "success", resMap["status"])
		assert.Equal(t, "Successfully update user!", resMap["message"])
	})
}
