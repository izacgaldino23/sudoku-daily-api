package bootstrap

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/middlewares"
	"sudoku-daily-api/src/infrastructure/logging"
)

func (c *Container) BuildRouters(app fiber.Router) {
	applyMiddlewares(app, c)

	// sudoku router
	sudokuGroup := app.Group("/sudoku")
	sudokuGroup.Get("/", c.Middlewares.Session, c.Middlewares.OptionalJWT, c.SudokuHandler.GetDailySudoku)
	sudokuGroup.Post("/generate", c.Middlewares.AuthOIDC, c.SudokuHandler.CreateSudoku)
	sudokuGroup.Post("/submit", c.Middlewares.Session, c.Middlewares.OptionalJWT, c.Middlewares.AuthMinimum, c.SudokuHandler.VerifySolution)
	sudokuGroup.Get("/me", c.Middlewares.RequireJWT, c.SudokuHandler.GetMyDailySudoku)

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
	leaderboardGroup.Post("/reset", c.Middlewares.AuthOIDC, c.LeaderboardHandler.ResetStrikes)
}

func applyMiddlewares(app fiber.Router, container *Container) {
	app.Use(recover.New(recover.Config{
		StackTraceHandler: func(c fiber.Ctx, e any) {
			stack := debug.Stack()
			logging.Log(c.Context()).Error().Msgf("Recovered from panic: %v\n%s", e, stack)

			_ = pkg.JsonError(c, pkg.ErrInternalServerError)
		},
		EnableStackTrace: true,
	}))

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

	app.Use(container.Middlewares.Timeout)

	app.Use(container.Middlewares.GlobalRateLimiter)
	app.Use(container.Middlewares.UserRateLimiter)

	app.Use(container.Middlewares.RequestID)
	app.Use(container.Middlewares.ResponseHeaders)
	app.Use(container.Middlewares.LogMiddleware)
}
