package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthLogout(t *testing.T) {
	setupTokens := func() (string, string) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		tokens, err := RegisterAndLoginUserWithTokens(app, "test@example.com", "testuser", "password123")
		assert.NoError(t, err)

		return tokens.AccessToken, tokens.RefreshToken
	}

	t.Run("valid logout", func(t *testing.T) {
		accessToken, refreshToken := setupTokens()
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		_, refreshToken := setupTokens()
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid access token", func(t *testing.T) {
		_, refreshToken := setupTokens()
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "invalid-token")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("missing refresh token", func(t *testing.T) {
		accessToken, _ := setupTokens()
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", accessToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
