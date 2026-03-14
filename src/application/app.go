package application

import (
	"sudoku-daily-api/src/application/bootstrap"
	"sudoku-daily-api/src/infrastructure/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func InitApp(app fiber.Router) error {
	container := &bootstrap.Container{}

	container.BuildInfrastructure()
	container.BuildRepositories()
	container.BuildServices()
	container.BuildUseCases()
	container.BuildHandlers()
	container.BuildMiddlewares()
	
	app.Use(container.Middlewares.ResponseHeaders)
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{"Origin",
			"Content-Type",
			"Accept",
			"X-Session-ID",
			"Authorization",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"X-Session-ID",
			"X-Request-ID",
		},
	}))

	app.Use(container.Middlewares.LogMiddleware)
	app.Use(container.Middlewares.RequestID)

	http.RegisterRoutes(
		app,
		container.SudokuHandler,
		container.AuthHandler,
		container.Middlewares.OptionalJWT,
		container.Middlewares.RequireJWT,
		container.Middlewares.Session,
		container.Middlewares.AuthMinimum,
	)

	return nil
}
