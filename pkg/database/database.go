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
		return fmt.Errorf("Error connecting to database: %w", err)
	}

	sqlDB.SetMaxIdleConns(configEnv.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(configEnv.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(configEnv.Database.MaxLifetime * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), configEnv.Database.Timeout*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping timeout: %w", err)
	}

	dbConnection.SqlConnection = sqlDB
	dbConnection.BunConnection = bun.NewDB(sqlDB, pgdialect.New())

	if configEnv.Debug {
		dbConnection.BunConnection.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return
}

func GetDB() DatabaseConnection {
	if dbConnection.SqlConnection == nil {
		_ = ConnectDB(config.GetConfig())
	}

	return dbConnection
}
