package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"sudoku-daily-api/pkg/config"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type DatabaseConnection struct {
	SqlConnection *sql.DB
	BunConnection *bun.DB
	Timeout       time.Duration
}

var (
	dbConnection DatabaseConnection
)

func ConnectDB(configEnv *config.Config) (err error) {
	dsn := configEnv.Database.DSNPostgres()

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	sqlDB.SetMaxIdleConns(configEnv.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(configEnv.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(configEnv.Database.MaxLifetime)

	dbConnection.SqlConnection = sqlDB
	dbConnection.BunConnection = bun.NewDB(sqlDB, pgdialect.New())

	if configEnv.Debug || configEnv.LogLevel == "debug" {
		dbConnection.BunConnection.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return pingWithRetry(sqlDB)
}

const pingTimeout = 5 * time.Second

func pingWithRetry(sqlDB *sql.DB) error {
	var err error
	for i := range 2 {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		err = sqlDB.PingContext(ctx)
		cancel()
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return fmt.Errorf("database not reachable after retries: %w", err)
}

func CloseDB() {
	if dbConnection.SqlConnection != nil {
		_ = dbConnection.SqlConnection.Close()
	}
}

func GetDB() DatabaseConnection {
	if dbConnection.SqlConnection == nil {
		_ = ConnectDB(config.GetConfig())
	}

	return dbConnection
}
