package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/infrastructure/http/sudoku"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSudokuGenerate(t *testing.T) {
	t.Run("generate daily sudoku", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/generate", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var sudokus []sudoku.SudokuResponse
		err = json.NewDecoder(resp.Body).Decode(&sudokus)
		assert.NoError(t, err)
		assert.NotEmpty(t, sudokus)

		for _, s := range sudokus {
			assert.NotEmpty(t, s.ID)
			assert.Greater(t, s.Size, 0)
			assert.NotEmpty(t, s.Board)
			assert.NotEmpty(t, s.Date)
		}
	})
}

func TestSudokuGetDaily(t *testing.T) {
	t.Cleanup(TruncateTables)
	app := SetupTestApp()

	req := httptest.NewRequest(http.MethodPost, "/api/sudoku/generate", nil)
	req.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(req)

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
}

func TestSudokuGetMyDailySolves(t *testing.T) {
	t.Run("get solves for today with data returns entries", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		err := SeedSudokus()
		assert.NoError(t, err)

		token, err := RegisterAndLoginUser(app, "user_me@example.com", "testuser", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		userID, err := GetUserIDByEmail("user_me@example.com")
		assert.NoError(t, err)

		err = SeedSolve(userID, sudokusIDs[0], 60)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var mySolves sudoku.MySolvesResponse
		err = json.NewDecoder(resp.Body).Decode(&mySolves)
		assert.NoError(t, err)
		if assert.Len(t, mySolves.Solves, 1) {
			assert.Equal(t, 60, mySolves.Solves[0].Duration)
		}
	})

	t.Run("get solves for today with no data returns empty list", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		token, err := RegisterAndLoginUser(app, "user_empty@example.com", "emptyuser", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var mySolves sudoku.MySolvesResponse
		err = json.NewDecoder(resp.Body).Decode(&mySolves)
		assert.NoError(t, err)
		assert.Empty(t, mySolves.Solves)
	})

	t.Run("get solves without auth returns 401", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("get solves with invalid token returns 401", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("get solves excludes yesterday's solves", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		err := SeedSudokus()
		assert.NoError(t, err)

		token, err := RegisterAndLoginUser(app, "user_yesterday@example.com", "yesterdayuser", "password123")
		assert.NoError(t, err)

		userID, err := GetUserIDByEmail("user_yesterday@example.com")
		assert.NoError(t, err)

		solve := SolveSeed{
			ID:        generateUUID(),
			UserID:    userID,
			SudokuID:  sudokusIDs[0],
			StartedAt: time.Now().Add(-26 * time.Hour),
			Duration:  90,
			Size:      9,
		}
		ctx := context.Background()
		_, err = database.GetDB().BunConnection.NewInsert().Model(&solve).Exec(ctx)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var mySolves sudoku.MySolvesResponse
		err = json.NewDecoder(resp.Body).Decode(&mySolves)
		assert.NoError(t, err)
		assert.Empty(t, mySolves.Solves)
	})
}

func TestSudokuSubmitWithoutLogin(t *testing.T) {
	t.Cleanup(TruncateTables)
	app := SetupTestApp()

	err := SeedSudokus()
	assert.NoError(t, err)

	solution, err := GetSudokuSolution(entities.BoardSize9)
	assert.NoError(t, err)

	dailyReq := httptest.NewRequest(http.MethodGet, "/api/sudoku?size=nine", nil)
	dailyReq.Header.Set("Content-Type", "application/json")
	dailyResp, err := app.Test(dailyReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, dailyResp.StatusCode)

	var sudokuResp sudoku.SudokuResponse
	err = json.NewDecoder(dailyResp.Body).Decode(&sudokuResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, sudokuResp.PlayToken)

	sessionID := dailyResp.Header.Get("X-Session-Id")
	assert.NotEmpty(t, sessionID)

	t.Run("submit valid solution", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"solution":   solution,
			"play_token": sudokuResp.PlayToken,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/submit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Session-Id", sessionID)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		respBody, _ := io.ReadAll(resp.Body)
		t.Logf("expected=%d, got=%d, body=%s", http.StatusOK, resp.StatusCode, string(respBody))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("submit with invalid session token", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"solution":   solution,
			"play_token": "invalid-token",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/submit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Session-Id", "invalid-token")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		respBody, _ := io.ReadAll(resp.Body)
		t.Logf("expected=%d, got=%d, body=%s", http.StatusUnauthorized, resp.StatusCode, string(respBody))
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("submit with missing session token", func(t *testing.T) {
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
		req.Header.Set("X-Session-Id", sessionID)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		respBody, _ := io.ReadAll(resp.Body)
		t.Logf("expected=%d, got=%d, body=%s", http.StatusBadRequest, resp.StatusCode, string(respBody))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
