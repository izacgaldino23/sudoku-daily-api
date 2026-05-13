package main

import (
	"os"
	"time"

	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/application"
	"sudoku-daily-api/src/infrastructure/http/middlewares"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "sudoku-daily-api/docs"
)

func init() {
	initLogger()

	err := config.Load()
	if err != nil {
		log.Logger.Fatal().Err(err)
	}

	c := config.GetConfig()
	err = database.ConnectDB(c)
	if err != nil {
		log.Logger.Fatal().Err(err)
	}
}

func main() {
	app := fiber.New()
	healthCheck(app)

	app.Get("/swagger/*", swaggo.HandlerDefault)
	app.Get("/metrics", middlewares.MetricsHandler())

	apiRouter := app.Group("/api")
	_ = application.InitApp(apiRouter)

	port := config.GetConfig().ApiPort
	log.Logger.Info().Msgf("Server running on port %v", port)

	err := app.Listen(port)
	if err != nil {
		log.Logger.Fatal().Err(err)
	}
}

func initLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	cfg := config.GetConfig()

	switch cfg.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		if cfg.Debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	}

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func healthCheck(app *fiber.App) {
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
