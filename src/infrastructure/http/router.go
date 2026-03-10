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
	optionalJWTMiddleware fiber.Handler,
	requireJWTMiddleware fiber.Handler,
	sessionMiddleware fiber.Handler,
	authMinimumMiddleware fiber.Handler,
) {
	registerSudokuRoutes(app.Group("/sudoku"), sudokuHandler, optionalJWTMiddleware, sessionMiddleware, authMinimumMiddleware)
	registerAuthRoutes(app.Group("/auth"), authHandler, optionalJWTMiddleware, requireJWTMiddleware)
}

func registerSudokuRoutes(
	api fiber.Router,
	sudokuHandler sudoku.ISudokuHandler,
	optionalJWTMiddleware fiber.Handler,
	sessionMiddleware fiber.Handler,
	authMinimumMiddleware fiber.Handler,
) {
	api.Get("/", sessionMiddleware, optionalJWTMiddleware, sudokuHandler.GetDailySudoku)
	api.Post("/generate", sudokuHandler.CreateSudoku)
	api.Post("/submit", sessionMiddleware, optionalJWTMiddleware, authMinimumMiddleware, sudokuHandler.VerifySolution)
}

func registerAuthRoutes(
	app fiber.Router,
	authHandler auth.AuthHandler,
	optionalJWTMiddleware fiber.Handler,
	requireJWTMiddleware fiber.Handler,
) {
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)

	private := app.Group("/", optionalJWTMiddleware, requireJWTMiddleware)

	private.Post("/refresh", authHandler.Refresh)
	private.Post("/logout", authHandler.Logout)
}
