package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"sudoku-daily-api/migrations"
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
)

func init() {
	err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading configs")
	}

	initLogger()

	if err = database.ConnectDB(config.GetConfig()); err != nil {
		log.Fatal().Err(err).Msg("Error connecting on database")
	}
}

func main() {
	if err := migrations.RunMigrations(config.GetConfig().Database.MigrationsPath); err != nil {
		log.Fatal().Err(err).Msg("Error running migrations")
	}

	log.Info().Msg("migrations executed successfully")
}

func initLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if config.GetConfig().Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
