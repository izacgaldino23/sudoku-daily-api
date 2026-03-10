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
	sessionMiddleware fiber.Handler,
) {
	registerSudokuRoutes(app.Group("/sudoku"), sudokuHandler, tokenMiddleware, sessionMiddleware)
	registerAuthRoutes(app.Group("/auth"), authHandler, tokenMiddleware, authMiddleware)
}

func registerSudokuRoutes(
	api fiber.Router,
	sudokuHandler sudoku.ISudokuHandler,
	tokenMiddleware fiber.Handler,
	sessionMiddleware fiber.Handler,
) {
	api.Get("/", sessionMiddleware, tokenMiddleware, sudokuHandler.GetDailySudoku)
	api.Post("/generate", sudokuHandler.CreateSudoku)
	api.Post("/submit", sudokuHandler.VerifySolution)
}

func registerAuthRoutes(
	app fiber.Router,
	authHandler auth.AuthHandler,
	tokenMiddleware fiber.Handler,
	authMiddleware fiber.Handler,
) {
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)

	private := app.Group("/", tokenMiddleware, authMiddleware)

	private.Post("/refresh", authHandler.Refresh)
	private.Post("/logout", authHandler.Logout)
}
