package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/tests/integration/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestAuthLogin(t *testing.T) {
	t.Run("valid login", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		_, err := testhelpers.RegisterAndLoginUser(app, "test@example.com", "testuser", "password123")
		assert.NoError(t, err)

		loginBody, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
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
		assert.NotEmpty(t, loginResp.RefreshToken)
		assert.Equal(t, "testuser", loginResp.UserName)
		assert.Equal(t, "test@example.com", loginResp.Email)
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		_, err := testhelpers.RegisterAndLoginUser(app, "test@example.com", "testuser", "password123")
		assert.NoError(t, err)

		loginBody, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("user not found", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

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
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		_, err := testhelpers.RegisterAndLoginUser(app, "test@example.com", "testuser", "password123")
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
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		_, err := testhelpers.RegisterAndLoginUser(app, "test@example.com", "testuser", "password123")
		assert.NoError(t, err)

		loginBody, _ := json.Marshal(map[string]string{
			"email": "test@example.com",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
