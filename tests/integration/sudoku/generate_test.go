package sudoku_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	"sudoku-daily-api/src/infrastructure/http/sudoku"
	"sudoku-daily-api/tests/integration/testhelpers"
)

func TestSudokuGenerate(t *testing.T) {
	t.Run("generate daily sudoku", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

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
