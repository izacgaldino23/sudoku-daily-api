package integration

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sudoku-daily-api/migrations"
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/application"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
)

func TestMain(m *testing.M) {
	setupTestEnvironment()
	defer teardownTestDB()
	os.Exit(m.Run())
}

func setupTestEnvironment() {
	os.Setenv("ENV", "test")

	os.Setenv("DATABASE_MIGRATIONS_PATH", "../../migrations/sql")

	os.Setenv("DATABASE_HOST", "127.0.0.1")
	os.Setenv("DATABASE_PORT", "5333")
	os.Setenv("DATABASE_USERNAME", "postgres")
	os.Setenv("DATABASE_PASSWORD", "12345")
	os.Setenv("DATABASE_NAME", "sudoku_test")
	os.Setenv("DATABASE_SSL_MODE", "disable")
	os.Setenv("API_PORT", "8081")

	memory := 64

	os.Setenv("AUTH_ITERATIONS", "3")
	os.Setenv("AUTH_MEMORY", strconv.Itoa(memory))
	os.Setenv("AUTH_PARALLELISM", "4")
	os.Setenv("AUTH_KEY_LEN", "32")
	os.Setenv("AUTH_SALT_LEN", "16")
	os.Setenv("AUTH_SECRET_KEY", "test-secret-key-for-integration-tests")
	os.Setenv("AUTH_ACCESS_TOKEN_DURATION", "15")
	os.Setenv("AUTH_REFRESH_TOKEN_DURATION", "60")

	err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	err = waitForDB(10, 2*time.Second)
	if err != nil {
		panic(fmt.Sprintf("failed to wait for database: %v", err))
	}

	database.GetDB().SqlConnection.SetConnMaxLifetime(5 * time.Minute)
	database.GetDB().SqlConnection.SetConnMaxIdleTime(2 * time.Minute)

	err = migrations.RunMigrations(config.GetConfig().Database.MigrationsPath)
	if err != nil {
		panic(fmt.Sprintf("failed to run migrations: %v", err))
	}
}

func waitForDB(maxRetries int, delay time.Duration) error {
	connected := false
	for i := 0; i < maxRetries; i++ {
		err := database.ConnectDB(config.GetConfig())
		if err == nil {
			connected = true
			break
		}

		time.Sleep(delay)
	}

	if connected {
		for i := 0; i < maxRetries; i++ {
			err := database.GetDB().SqlConnection.Ping()
			if err == nil {
				return nil
			}

			time.Sleep(delay)
		}
	}

	return fmt.Errorf("database not available after %d retries", maxRetries)
}

func teardownTestDB() {
	dbConn := database.GetDB()
	if dbConn.SqlConnection != nil {
		dbConn.SqlConnection.Close()
	}
}

func TruncateTables(t *testing.T) {
	dbConn := database.GetDB()
	if dbConn.BunConnection == nil {
		return
	}

	tables := []string{
		`"solves"`,
		`"refresh_tokens"`,
		`"users"`,
		`"sudokus"`,
	}

	ctx := context.Background()
	for _, table := range tables {
		_, err := dbConn.BunConnection.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			t.Fatalf("warning: failed to truncate table %s: %v\n", table, err)
		}
	}
}

func SetupTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
	})

	api := app.Group("/api")

	_ = application.InitApp(api)

	return app
}
