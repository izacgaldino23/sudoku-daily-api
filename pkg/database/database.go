package database

import (
	"database/sql"
	"fmt"
	"sudoku-daily-api/pkg/config"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type DatabaseConnection struct {
	SqlConnection *sql.DB
	BunConnection *bun.DB
}

var (
	dbConnection DatabaseConnection
)

func ConnectDB(configEnv *config.Config) (err error) {
	dsn := configEnv.Database.DSNPostgres()
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %w", err)
	}

	dbConnection.SqlConnection = sqlDB
	dbConnection.BunConnection = bun.NewDB(sqlDB, pgdialect.New())

	return
}

func GetDB() DatabaseConnection {
	if dbConnection.SqlConnection == nil {
		_ = ConnectDB(config.GetConfig())
	}

	return dbConnection
}
