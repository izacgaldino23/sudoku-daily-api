package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/stretchr/testify/assert"
)

func TestAuthLogin(t *testing.T) {
	t.Run("valid login", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		userData, err := helpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)

		loginBody, _ := json.Marshal(map[string]string{
			"email":    userData.Email,
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp auth.LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, loginResp.AccessToken)
		assert.Equal(t, userData.Username, loginResp.UserName)
		assert.Equal(t, userData.Email, loginResp.Email)

		cookies := resp.Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, "refresh_token", cookies[0].Name)
		assert.NotEmpty(t, cookies[0].Value)
		assert.True(t, cookies[0].HttpOnly)
		assert.Equal(t, "/", cookies[0].Path)
	})

	t.Run("valid login sets cookie with correct expiry", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		userData, err := helpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)

		loginBody, _ := json.Marshal(map[string]string{
			"email":    userData.Email,
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cookies := resp.Cookies()
		assert.Len(t, cookies, 1)
		cookie := cookies[0]

		refreshTokenDuration := config.GetConfig().Auth.RefreshTokenDuration
		expectedMaxAge := refreshTokenDuration
		assert.Equal(t, expectedMaxAge, cookie.MaxAge)
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		userData, err := helpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)

		loginBody, _ := json.Marshal(map[string]string{
			"email":    userData.Email,
			"password": "wrongpassword",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("user not found", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		loginBody, _ := json.Marshal(map[string]string{
			"email":    "nonexistent@example.com",
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("missing email", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		_, err := helpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)

		loginBody, _ := json.Marshal(map[string]string{
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing password", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		email := helpers.GenerateUniqueEmail("test")
		_, err := helpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)

		loginBody, _ := json.Marshal(map[string]string{
			"email": email,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}