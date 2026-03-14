package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/sudoku"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSudokuGenerate(t *testing.T) {
	TruncateTables(t)

	app := SetupTestApp()

	tests := []struct {
		name       string
		wantStatus int
	}{
		{
			name:       "generate daily sudoku",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TruncateTables(t)

			req := httptest.NewRequest(http.MethodPost, "/api/sudoku/generate", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, fiber.TestConfig{
				Timeout: 0,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var sudokus []sudoku.SudokuResponse
				err := json.NewDecoder(resp.Body).Decode(&sudokus)
				assert.NoError(t, err)
				assert.NotEmpty(t, sudokus)

				for _, s := range sudokus {
					assert.NotEmpty(t, s.ID)
					assert.Greater(t, s.Size, 0)
					assert.NotEmpty(t, s.Board)
					assert.NotEmpty(t, s.Date)
				}
			} else {
				result := pkg.Error{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)

				if tt.wantStatus != http.StatusOK {
					assert.NotEmpty(t, result.Message)
				} else {
					assert.Empty(t, result.Message)
				}
			}
		})
	}
}

func TestSudokuGetDaily(t *testing.T) {
	TruncateTables(t)

	app := SetupTestApp()

	req := httptest.NewRequest(http.MethodPost, "/api/sudoku/generate", nil)
	req.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(req)

	tests := []struct {
		name       string
		session    string
		auth       string
		size       string
		wantStatus int
	}{
		{
			name:       "get daily sudoku without login",
			session:    "",
			auth:       "",
			size:       "nine",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get daily sudoku with session header",
			session:    uuid.NewString(),
			auth:       "",
			size:       "nine",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get daily sudoku with invalid size",
			session:    "",
			auth:       "",
			size:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "get daily sudoku with missing size",
			session:    "",
			auth:       "",
			size:       "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "get daily sudoku with size four",
			session:    "",
			auth:       "",
			size:       "four",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get daily sudoku with size six",
			session:    "",
			auth:       "",
			size:       "six",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/sudoku"
			if tt.size != "" {
				url += "?size=" + tt.size
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")

			if tt.session != "" {
				req.Header.Set("X-Session-Id", tt.session)
			}
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}

			resp, err := app.Test(req, fiber.TestConfig{
				Timeout: 0,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			sessionID := resp.Header.Get("X-Session-Id")

			if resp.StatusCode == http.StatusOK {
				var sudokuResp sudoku.SudokuResponse
				err := json.NewDecoder(resp.Body).Decode(&sudokuResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, sudokuResp.ID)
				assert.NotEmpty(t, sudokuResp.PlayToken)
				assert.NotEmpty(t, sudokuResp.Board)
				assert.NotEmpty(t, sudokuResp.Date)
				assert.NotEmpty(t, sessionID)
			} else {
				var errResp pkg.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				assert.NoError(t, err)
				if tt.wantStatus != http.StatusOK {
					assert.NotEmpty(t, errResp.Message)
				}
			}
		})
	}
}

func TestSudokuSubmitWithoutLogin(t *testing.T) {
	TruncateTables(t)

	app := SetupTestApp()

	err := SeedSudokus()
	assert.NoError(t, err)

	solution, err := GetSudokuSolution(9)
	assert.NoError(t, err)
	t.Logf("Solution from DB: %v", solution)
	t.Logf("Solution type: %T", solution)

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

	tests := []struct {
		name       string
		body       map[string]interface{}
		header     string
		wantStatus int
	}{
		{
			name: "submit valid solution",
			body: map[string]interface{}{
				"solution":   solution,
				"play_token": sudokuResp.PlayToken,
			},
			header:     sessionID,
			wantStatus: http.StatusOK,
		},
		{
			name: "submit with invalid session token",
			body: map[string]interface{}{
				"solution":   solution,
				"play_token": "invalid-token",
			},
			header:     "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "submit with missing session token",
			body: map[string]interface{}{
				"solution": solution,
				"play_token": sudokuResp.PlayToken,
			},
			header:     "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "submit with missing solution",
			body: map[string]interface{}{
				"play_token": sudokuResp.PlayToken,
			},
			header:     sessionID,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/sudoku/submit", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.header != "" {
				req.Header.Set("X-Session-Id", tt.header)
			}

			resp, err := app.Test(req, fiber.TestConfig{
				Timeout: 0,
			})
			assert.NoError(t, err)

			respBody, _ := io.ReadAll(resp.Body)
			t.Logf("Test '%s': expected=%d, got=%d, body=%s", tt.name, tt.wantStatus, resp.StatusCode, string(respBody))
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
