package migrations

import (
	"errors"
	"fmt"

	"sudoku-daily-api/pkg/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

type migrateLogger struct{}

func (m *migrateLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf("migrate: "+format, v...)
}

func (m *migrateLogger) Verbose() bool {
	return true
}

func RunMigrations(migrationsPath string) error {
	driver, err := postgres.WithInstance(database.GetDB().SqlConnection, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migrate: %w", err)
	}

	m.Log = &migrateLogger{}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Info().Msg("migrations executed successfully")
	return nil
}
