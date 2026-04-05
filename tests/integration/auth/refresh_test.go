package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/tests/integration/testhelpers"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthRefresh(t *testing.T) {
	setupTokens := func() (string, string) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		email := testhelpers.GenerateUniqueEmail("test")
		tokens, err := testhelpers.RegisterAndLoginUserWithTokens(app, email, "testuser", "password123")
		assert.NoError(t, err)

		return tokens.AccessToken, tokens.RefreshToken
	}

	t.Run("valid refresh", func(t *testing.T) {
		accessToken, refreshToken := setupTokens()
		app := testhelpers.SetupTestApp()

		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", accessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var refreshResp auth.RefreshTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&refreshResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, refreshResp.AccessToken)
	})

	t.Run("missing refresh token", func(t *testing.T) {
		accessToken, _ := setupTokens()
		app := testhelpers.SetupTestApp()

		body, _ := json.Marshal(map[string]string{})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", accessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var result pkg.Error
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		accessToken, _ := setupTokens()
		app := testhelpers.SetupTestApp()

		body, _ := json.Marshal(map[string]string{"refresh_token": "invalid-refresh-token"})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", accessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var result pkg.Error
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.NotEmpty(t, result.Message)
	})
}
