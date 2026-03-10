package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/sudoku"
	"testing"

	"github.com/gofiber/fiber/v3"
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
			session:    "some-session-token",
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
				req.Header.Set("session", tt.session)
			}
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}

			resp, err := app.Test(req, fiber.TestConfig{
				Timeout: 0,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var sudokuResp sudoku.SudokuResponse
				err := json.NewDecoder(resp.Body).Decode(&sudokuResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, sudokuResp.ID)
				assert.NotEmpty(t, sudokuResp.SessionToken)
				assert.NotEmpty(t, sudokuResp.Board)
				assert.NotEmpty(t, sudokuResp.Date)
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
