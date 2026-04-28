package sudoku_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/infrastructure/http/sudoku"
	"sudoku-daily-api/tests/integration/helpers"
)

func boardSizeToName(size entities.BoardSize) string {
	switch size {
	case entities.BoardSize4:
		return "four"
	case entities.BoardSize6:
		return "six"
	case entities.BoardSize9:
		return "nine"
	default:
		return "unknown"
	}
}

func TestSudokuGenerate(t *testing.T) {
	t.Run("generate daily sudoku without size", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/generate", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("generate daily sudoku with invalid size", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/sudoku/generate?size=invalid", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	for boardSize := range entities.BoardSizes {
		t.Run("generate daily sudoku with size "+boardSizeToName(boardSize), func(t *testing.T) {
			t.Cleanup(helpers.TruncateTables)
			app := helpers.SetupTestApp()

			sizeName := boardSizeToName(boardSize)
			req := httptest.NewRequest(http.MethodPost, "/api/sudoku/generate?size="+sizeName, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var sudoku sudoku.SudokuResponse
			err = json.NewDecoder(resp.Body).Decode(&sudoku)
			assert.NoError(t, err)
			assert.NotEmpty(t, sudoku.ID)
			assert.Equal(t, int(boardSize), sudoku.Size)
			assert.NotEmpty(t, sudoku.Board)
			assert.NotEmpty(t, sudoku.Date)
		})
	}
}
