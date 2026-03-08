package migrations

import (
	"errors"
	"fmt"
	"sudoku-daily-api/pkg/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(migrationsPath string) error {
	driver, err := postgres.WithInstance(database.GetDB().SqlConnection, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Error creating driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
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
