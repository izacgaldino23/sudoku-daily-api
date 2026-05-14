package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"sudoku-daily-api/migrations"
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

	if c.Database.MigrationsEnabled {
		if err = migrations.RunMigrations(c.Database.MigrationsPath); err != nil {
			log.Logger.Fatal().Err(err).Msg("Error running migrations")
		}
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
	cfg := config.GetConfig()

	go func() {
		log.Logger.Info().Msgf("Server running on port %v", port)
		if err := app.Listen(port); err != nil {
			log.Logger.Fatal().Err(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Logger.Info().Msg("Shutting down server...")

	if err := app.ShutdownWithTimeout(cfg.Limits.ShutdownTimeout); err != nil {
		log.Logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	database.CloseDB()

	log.Logger.Info().Msg("Server exited")
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
