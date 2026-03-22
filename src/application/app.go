package application

import (
	"sudoku-daily-api/src/application/bootstrap"

	"github.com/gofiber/fiber/v3"
)

func InitApp(app fiber.Router) error {
	container := &bootstrap.Container{}

	container.BuildInfrastructure()
	container.BuildRepositories()
	container.BuildServices()
	container.BuildUseCases()
	container.BuildHandlers()
	container.BuildMiddlewares()
	container.BuildRouters(app)

	return nil
}
