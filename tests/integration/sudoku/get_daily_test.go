package sudoku_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/sudoku"
	"sudoku-daily-api/tests/integration/helpers"
)

func TestSudokuGetDaily(t *testing.T) {
	t.Cleanup(helpers.TruncateTables)
	app := helpers.SetupTestApp()

	req := httptest.NewRequest(http.MethodPost, "/api/sudoku/generate", nil)
	req.Header.Set("Content-Type", "application/json")
	_, err := app.Test(req, fiber.TestConfig{Timeout: 0})
	assert.NoError(t, err)

	t.Run("get daily sudoku without login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/sudoku?size=nine", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, _ := io.ReadAll(resp.Body)
		t.Logf("status=%d, body=%s", resp.StatusCode, string(respBody))

		var sudokuResp sudoku.SudokuResponse
		err = json.Unmarshal(respBody, &sudokuResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, sudokuResp.ID)
		assert.NotEmpty(t, sudokuResp.PlayToken)
		assert.NotEmpty(t, sudokuResp.Board)
		assert.NotEmpty(t, sudokuResp.Date)
	})

	t.Run("get daily sudoku with session header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/sudoku?size=nine", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Session-Id", uuid.NewString())

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("get daily sudoku with invalid size", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/sudoku?size=invalid", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp pkg.Error
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, errResp.Message)
	})

	t.Run("get daily sudoku with missing size", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/sudoku", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get daily sudoku with size four", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/sudoku?size=four", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("get daily sudoku with size six", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/sudoku?size=six", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("get daily sudoku with size four and user logged in", func(t *testing.T) {
		userData, err := helpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, userData.AccessToken)

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku?size=four", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userData.AccessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var sudokuResp sudoku.SudokuResponse
		err = json.Unmarshal(respBody, &sudokuResp)
		assert.NoError(t, err)

		assert.NotEmpty(t, sudokuResp.ID)
		assert.NotEmpty(t, sudokuResp.Size)
		assert.NotEmpty(t, sudokuResp.Board)
		assert.NotEmpty(t, sudokuResp.Date)
		assert.NotEmpty(t, sudokuResp.PlayToken)
		assert.Empty(t, sudokuResp.SessionID)
	})
}
