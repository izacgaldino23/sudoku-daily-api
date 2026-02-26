package http

import (
	"sudoku-daily-api/src/adapters/inbound/http/sudoku"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(api fiber.Router, sudokuHandler sudoku.ISudokuHandler) {
	api.Get("/sudoku/:size", sudokuHandler.GetDailySudoku)
}
