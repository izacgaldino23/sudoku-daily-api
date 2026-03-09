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
	tokenMiddleware fiber.Handler,
	authMiddleware fiber.Handler,
) {
	registerSudokuRoutes(app, sudokuHandler, tokenMiddleware)
	registerAuthRoutes(app, authHandler, tokenMiddleware, authMiddleware)
}

func registerSudokuRoutes(api fiber.Router, sudokuHandler sudoku.ISudokuHandler, tokenMiddleware fiber.Handler) {
	api.Get("/sudoku", sudokuHandler.GetDailySudoku)
	api.Post("/sudoku/generate", sudokuHandler.CreateSudoku)
}

func registerAuthRoutes(
	app fiber.Router,
	authHandler auth.AuthHandler,
	tokenMiddleware fiber.Handler,
	authMiddleware fiber.Handler,
) {
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)

	private := app.Group("/auth", tokenMiddleware, authMiddleware)

	private.Post("/refresh", authHandler.Refresh)
	private.Post("/logout", authHandler.Logout)
}
