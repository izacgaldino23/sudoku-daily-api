package database

import (
	"database/sql"
	"fmt"
	"sudoku-daily-api/pkg/config"

	_ "github.com/lib/pq"
)

type DatabaseConnection struct {
	SqlConnection  *sql.DB
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

	return
}

func GetDB() *sql.DB {
	return dbConnection.SqlConnection
}