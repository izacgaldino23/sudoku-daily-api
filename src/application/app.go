package application

import (
	"sudoku-daily-api/src/application/bootstrap"
	"sudoku-daily-api/src/infrastructure/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func InitApp(app fiber.Router) error {

	app.Use(recover.New())

	container := &bootstrap.Container{}

	container.BuildInfrastructure()
	container.BuildRepositories()
	container.BuildServices()
	container.BuildUseCases()
	container.BuildHandlers()
	container.BuildMiddlewares()

	app.Use(container.LogMiddleware)

	http.RegisterRoutes(
		app,
		container.SudokuHandler,
		container.AuthHandler,
		container.OptionalJWT,
		container.RequireJWT,
		container.Session,
		container.AuthMinimum,
	)

	return nil
}
