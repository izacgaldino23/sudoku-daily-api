package application

import (
	"sudoku-daily-api/src/application/bootstrap"
	"sudoku-daily-api/src/infrastructure/http"
	"sudoku-daily-api/src/infrastructure/http/middlewares"

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

	addMiddlewares(app, container)

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

func addMiddlewares(app fiber.Router, container *bootstrap.Container) {
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
