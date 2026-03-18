package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthRegister(t *testing.T) {
	t.Run("valid registration", func(t *testing.T) {
		TruncateTables(t)
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
			"username": "testuser",
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("missing email", func(t *testing.T) {
		TruncateTables(t)
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{
			"username": "testuser",
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid email format", func(t *testing.T) {
		TruncateTables(t)
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{
			"email":    "invalid-email",
			"username": "testuser",
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing username", func(t *testing.T) {
		TruncateTables(t)
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("username too short", func(t *testing.T) {
		TruncateTables(t)
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
			"username": "ab",
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing password", func(t *testing.T) {
		TruncateTables(t)
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
			"username": "testuser",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("password too short", func(t *testing.T) {
		TruncateTables(t)
		app := SetupTestApp()

		body, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
			"username": "testuser",
			"password": "123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
