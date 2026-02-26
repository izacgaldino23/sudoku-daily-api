package migrations

import (
	"errors"
	"fmt"
	"sudoku-daily-api/pkg/database"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() error {
	driver, err := postgres.WithInstance(database.GetDB(), &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Error creating driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/sql",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Error creating migrate: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("Error running migrations: %w", err)
	}

	return nil
}
