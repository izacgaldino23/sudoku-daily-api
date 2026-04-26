package sudoku_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/infrastructure/http/sudoku"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestSudokuSubmitWithLogin(t *testing.T) {
	t.Cleanup(helpers.TruncateTables)
	app := helpers.SetupTestApp()

	err := helpers.SeedSudokus()
	assert.NoError(t, err)

	solution, err := helpers.GetSudokuSolution(entities.BoardSize9)
	assert.NoError(t, err)

	userData, err := helpers.RegisterAndLoginUser(app, "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, userData.AccessToken)

	userID, err := helpers.GetUserIDByEmail(userData.Email)
	assert.NoError(t, err)

	dailyReq := httptest.NewRequest(http.MethodGet, "/api/sudoku?size=nine", nil)
	dailyReq.Header.Set("Content-Type", "application/json")
	dailyReq.Header.Set("Authorization", "Bearer "+userData.AccessToken)
	dailyResp, err := app.Test(dailyReq, fiber.TestConfig{Timeout: 0})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, dailyResp.StatusCode)

	var sudokuResp sudoku.SudokuResponse
	err = json.NewDecoder(dailyResp.Body).Decode(&sudokuResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, sudokuResp.PlayToken)

	t.Run("submit valid solution", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"solution":   solution,
			"play_token": sudokuResp.PlayToken,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/submit", bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+userData.AccessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		respBody, _ := io.ReadAll(resp.Body)
		t.Logf("expected=%d, got=%d, body=%s", http.StatusOK, resp.StatusCode, string(respBody))
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		solves, streaks, err := helpers.GetUserStats(userID)
		assert.NoError(t, err)

		assert.Equal(t, 1, solves)
		assert.Equal(t, 1, streaks)
	})

	t.Run("submit with invalid authorization token", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"solution":   solution,
			"play_token": "invalid-token",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/submit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "invalid-token")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		respBody, _ := io.ReadAll(resp.Body)
		t.Logf("expected=%d, got=%d, body=%s", http.StatusUnauthorized, resp.StatusCode, string(respBody))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("submit with missing authorization token", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"solution":   solution,
			"play_token": sudokuResp.PlayToken,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/submit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		respBody, _ := io.ReadAll(resp.Body)
		t.Logf("expected=%d, got=%d, body=%s", http.StatusUnauthorized, resp.StatusCode, string(respBody))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("submit with missing solution", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"play_token": sudokuResp.PlayToken,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/submit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userData.AccessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		respBody, _ := io.ReadAll(resp.Body)
		t.Logf("expected=%d, got=%d, body=%s", http.StatusBadRequest, resp.StatusCode, string(respBody))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}