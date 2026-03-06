package http

import (
	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/src/infrastructure/http/sudoku"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(
	app fiber.Router,
	sudokuHandler sudoku.ISudokuHandler,
	authHandler auth.AuthHandler,
) {
	registerSudokuRoutes(app, sudokuHandler)
	registerAuthRoutes(app, authHandler)
}

func registerSudokuRoutes(api fiber.Router, sudokuHandler sudoku.ISudokuHandler) {
	api.Get("/sudoku", sudokuHandler.GetDailySudoku)
	api.Post("/sudoku/generate", sudokuHandler.CreateSudoku)
}

func registerAuthRoutes(app fiber.Router, authHandler auth.AuthHandler) {
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/refresh", authHandler.Refresh)
}
