package bootstrap

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"sudoku-daily-api/src/infrastructure/http/middlewares"
)

func (c *Container) BuildRouters(app fiber.Router) {
	addMiddlewares(app, c)

	// sudoku router
	sudokuGroup := app.Group("/sudoku")
	sudokuGroup.Get("/", c.Middlewares.Session, c.Middlewares.OptionalJWT, c.SudokuHandler.GetDailySudoku)
	sudokuGroup.Post("/generate", c.SudokuHandler.CreateSudoku)
	sudokuGroup.Post("/submit", c.Middlewares.Session, c.Middlewares.OptionalJWT, c.Middlewares.AuthMinimum, c.SudokuHandler.VerifySolution)

	// auth router
	authGroup := app.Group("/auth")
	authGroup.Post("/register", c.AuthHandler.Register)
	authGroup.Post("/login", c.AuthHandler.Login)
	authGroup.Post("/refresh", c.AuthHandler.Refresh)

	private := authGroup.Group("/", c.Middlewares.RequireJWT)

	private.Post("/logout", c.AuthHandler.Logout)
	private.Get("/resume", c.AuthHandler.Resume)

	// leaderboard router
	leaderboardGroup := app.Group("/leaderboard")
	leaderboardGroup.Get("/", c.LeaderboardHandler.GetLeaderboard)
}

func addMiddlewares(app fiber.Router, container *Container) {
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			middlewares.XSessionIdHeader,
			middlewares.XRequestIDHeader,
		},
		ExposeHeaders: []string{
			middlewares.XSessionIdHeader,
			middlewares.XRequestIDHeader,
		},
	}))

	app.Use(container.Middlewares.RequestID)
	app.Use(container.Middlewares.ResponseHeaders)
	app.Use(container.Middlewares.LogMiddleware)
}
