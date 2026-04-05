package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"sudoku-daily-api/migrations"
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/application"
	"sudoku-daily-api/src/application/bootstrap"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

var Container *bootstrap.Container

var emailCounter atomic.Int64

var SudokusIDs = []string{
	"00000000-0000-0000-0000-000000000001",
	"00000000-0000-0000-0000-000000000002",
	"00000000-0000-0000-0000-000000000003",
}

var (
	setupOnce    sync.Once
	teardownOnce sync.Once
	appMutex     sync.Mutex
	truncateMu   sync.Mutex
)

// SetupTestEnvironment initializes the test environment (config, DB, migrations).
// Safe to call from multiple TestMain functions - only runs once.
func SetupTestEnvironment() {
	setupOnce.Do(func() {
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_MIGRATIONS_PATH", "../../../migrations/sql")
		os.Setenv("DATABASE_HOST", "127.0.0.1")
		os.Setenv("DATABASE_PORT", "5333")
		os.Setenv("DATABASE_USERNAME", "postgres")
		os.Setenv("DATABASE_PASSWORD", "12345")
		os.Setenv("DATABASE_NAME", "sudoku_test")
		os.Setenv("DATABASE_SSL_MODE", "disable")
		os.Setenv("API_PORT", "8081")
		os.Setenv("DEBUG", "false")

		memory := 64
		os.Setenv("AUTH_ITERATIONS", "3")
		os.Setenv("AUTH_MEMORY", strconv.Itoa(memory))
		os.Setenv("AUTH_PARALLELISM", "4")
		os.Setenv("AUTH_KEY_LEN", "32")
		os.Setenv("AUTH_SALT_LEN", "16")
		os.Setenv("AUTH_SECRET_KEY", "test-secret-key-for-integration-tests")
		os.Setenv("AUTH_ACCESS_TOKEN_DURATION", "15")
		os.Setenv("AUTH_REFRESH_TOKEN_DURATION", "60")
		os.Setenv("LIMITS_MAX_REQUEST_COUNT_GLOBAL", "1000")
		os.Setenv("LIMITS_MAX_REQUEST_COUNT_USER", "100")

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
	})
}

// TeardownTestDB closes the database connection.
// Safe to call from multiple TestMain functions - only runs once.
func TeardownTestDB() {
	teardownOnce.Do(func() {
		dbConn := database.GetDB()
		if dbConn.SqlConnection != nil {
			dbConn.SqlConnection.Close()
		}
	})
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

func TruncateTables() {
	truncateMu.Lock()
	defer truncateMu.Unlock()

	dbConn := database.GetDB()
	if dbConn.BunConnection == nil {
		return
	}

	tables := []string{
		`"solves"`,
		`"refresh_tokens"`,
		`"users"`,
		`"sudokus"`,
		`"user_stats"`,
	}

	ctx := context.Background()
	for _, table := range tables {
		_, err := dbConn.BunConnection.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			log.Printf("warning: failed to truncate table %s: %v\n", table, err)
		}
	}
}

func SetupTestApp() *fiber.App {
	appMutex.Lock()
	defer appMutex.Unlock()

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
	})

	api := app.Group("/api")
	Container = application.InitApp(api)

	return app
}

// UserData holds user access and refresh tokens and data
type UserData struct {
	Email        string
	Username     string
	AccessToken  string
	RefreshToken string
}

// RegisterAndLoginUser registers a new user and returns their access token
func RegisterAndLoginUser(app *fiber.App, password string) (UserData, error) {
	email := GenerateUniqueEmail("test_mail")
	username := GenerateUniqueUsername("test_username")
	userData, err := RegisterAndLoginUserWithTokens(app, email, username, password)
	if err != nil {
		return UserData{}, err
	}
	return userData, nil
}

// RegisterAndLoginUserWithTokens registers a new user and returns both access and refresh tokens
func RegisterAndLoginUserWithTokens(app *fiber.App, email, username, password string) (UserData, error) {
	registerBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"username": username,
		"password": password,
	})

	registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	_, err := app.Test(registerReq)
	if err != nil {
		return UserData{}, err
	}

	creds, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(creds))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := app.Test(loginReq)
	if err != nil {
		return UserData{}, err
	}

	var loginRespBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(loginResp.Body).Decode(&loginRespBody); err != nil {
		return UserData{}, err
	}

	return UserData{
		AccessToken:  loginRespBody.AccessToken,
		RefreshToken: loginRespBody.RefreshToken,
		Email:        email,
		Username:     username,
	}, nil
}

type SudokuSeed struct {
	bun.BaseModel `bun:"table:sudokus"`

	ID         string    `bun:"id,pk"`
	Size       int       `bun:",notnull"`
	Difficulty string    `bun:",notnull"`
	Board      []byte    `bun:"type:,notnull"`
	Solution   []byte    `bun:"type:,notnull"`
	Date       time.Time `bun:"type:date,notnull"`
}

func SeedSudokus() error {
	db := database.GetDB().BunConnection
	ctx := context.Background()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	sudokus := []SudokuSeed{
		{
			ID:         SudokusIDs[0],
			Size:       9,
			Difficulty: "easy",
			Board:      []byte{0, 0, 0, 0, 9, 4, 0, 3, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 1, 8, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 2, 3, 0, 0, 0, 0, 0, 5, 0, 0, 0, 8, 0, 7, 0, 0, 4, 0, 0, 0, 0, 0, 0, 3, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 7, 4, 0, 0, 0, 0, 0, 0, 0},
			Solution:   []byte{5, 8, 7, 6, 9, 4, 1, 3, 2, 6, 1, 4, 3, 2, 7, 9, 5, 8, 9, 2, 3, 8, 1, 5, 4, 7, 6, 3, 7, 8, 2, 4, 9, 6, 1, 5, 4, 5, 9, 7, 3, 6, 8, 2, 1, 1, 4, 6, 5, 7, 3, 2, 8, 9, 2, 3, 1, 9, 6, 8, 5, 7, 4, 7, 9, 5, 1, 3, 6, 8, 2, 4, 8, 6, 2, 4, 5, 1, 3, 9, 7, 1, 5, 3, 7, 2, 9, 4, 6, 8},
			Date:       today,
		},
		{
			ID:         SudokusIDs[1],
			Size:       4,
			Difficulty: "easy",
			Board:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Solution:   []byte{1, 2, 3, 4, 3, 4, 1, 2, 2, 1, 4, 3, 4, 3, 2, 1},
			Date:       today,
		},
		{
			ID:         SudokusIDs[2],
			Size:       6,
			Difficulty: "easy",
			Board:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Solution:   []byte{1, 2, 3, 4, 5, 6, 4, 5, 6, 1, 2, 3, 2, 3, 4, 5, 6, 1, 5, 6, 1, 2, 3, 4, 3, 4, 5, 6, 1, 2, 6, 1, 2, 3, 4, 5},
			Date:       today,
		},
	}

	for _, s := range sudokus {
		_, err := db.NewInsert().Model(&s).On("CONFLICT DO NOTHING").Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetSudokuSolution(size entities.BoardSize) ([][]int, error) {
	db := database.GetDB().BunConnection
	ctx := context.Background()

	var sudoku SudokuSeed
	err := db.NewSelect().Model(&sudoku).Where("size = ?", size).Order("date DESC").Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return bytesToMatrix(sudoku.Solution), nil
}

func bytesToMatrix(data []byte) [][]int {
	size := int(math.Sqrt(float64(len(data))))

	matrix := make([][]int, size)
	for i := 0; i < size; i++ {
		row := make([]int, size)
		for j := 0; j < size; j++ {
			row[j] = int(data[i*size+j])
		}
		matrix[i] = row
	}

	return matrix
}

type SolveSeed struct {
	bun.BaseModel `bun:"table:solves"`

	ID        string    `bun:"id,pk"`
	UserID    string    `bun:"user_id,notnull"`
	SudokuID  string    `bun:"sudoku_id,notnull"`
	StartedAt time.Time `bun:"type:timestamp,notnull"`
	Duration  int       `bun:",notnull"`
	Size      int       `bun:",notnull"`
	CreatedAt time.Time `bun:"type:timestamp,notnull,default:current_timestamp"`
}

func SeedSolve(userID, sudokuID string, duration int) error {
	db := database.GetDB().BunConnection
	ctx := context.Background()

	size := 9
	switch sudokuID {
	case SudokusIDs[1]:
		size = 4
	case SudokusIDs[2]:
		size = 6
	}

	solve := SolveSeed{
		ID:        GenerateUUID(),
		UserID:    userID,
		SudokuID:  sudokuID,
		StartedAt: time.Now().Add(-time.Duration(duration) * time.Second),
		Duration:  duration,
		Size:      size,
	}

	_, err := db.NewInsert().Model(&solve).Exec(ctx)
	return err
}

func SeedSolves(userID string) error {
	solves := []SolveSeed{
		{ID: GenerateUUID(), UserID: userID, SudokuID: SudokusIDs[0], StartedAt: time.Now().Add(-60 * time.Second), Duration: 60, Size: 9},
		{ID: GenerateUUID(), UserID: userID, SudokuID: SudokusIDs[0], StartedAt: time.Now().Add(-120 * time.Second), Duration: 120, Size: 9},
		{ID: GenerateUUID(), UserID: userID, SudokuID: SudokusIDs[1], StartedAt: time.Now().Add(-30 * time.Second), Duration: 30, Size: 4},
		{ID: GenerateUUID(), UserID: userID, SudokuID: SudokusIDs[0], StartedAt: time.Now().Add(-24 * time.Hour), Duration: 90, Size: 9},
		{ID: GenerateUUID(), UserID: userID, SudokuID: SudokusIDs[2], StartedAt: time.Now().Add(-25 * time.Hour), Duration: 45, Size: 4},
	}

	ctx := context.Background()
	solveDate := time.Now().Truncate(24 * time.Hour)
	for _, s := range solves {
		_, err := database.GetDB().BunConnection.NewInsert().Model(&s).Exec(context.Background())
		if err != nil {
			return err
		}

		err = Container.UserStatsSolveAddStrike.Execute(ctx, vo.UUID(userID), solveDate)
		if err != nil {
			return fmt.Errorf("failed to add strike: %w", err)
		}
	}

	return nil
}

func GetUserIDByEmail(email string) (string, error) {
	db := database.GetDB().BunConnection
	ctx := context.Background()

	var user struct {
		ID string `bun:"id"`
	}
	err := db.NewSelect().Model(&user).Table("users").Where("email = ?", email).Scan(ctx)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

type UserSeed struct {
	bun.BaseModel `bun:"table:users"`

	ID            string `bun:"id,pk"`
	Email         string `bun:",notnull,unique"`
	Username      string `bun:",notnull,unique"`
	PasswordHash  string `bun:",notnull"`
	Provider      string `bun:",notnull"`
	EmailVerified bool   `bun:",notnull,default:false"`
	CreatedAt     string `bun:"type:timestamp,notnull,default:current_timestamp"`
}

func SeedUser(email, username, passwordHash string) error {
	db := database.GetDB().BunConnection
	ctx := context.Background()

	user := UserSeed{
		ID:            GenerateUUID(),
		Email:         email,
		Username:      username,
		PasswordHash:  passwordHash,
		Provider:      "email",
		EmailVerified: false,
	}

	_, err := db.NewInsert().Model(&user).Exec(ctx)
	return err
}

func GenerateUUID() string {
	return vo.NewUUID().String()
}

// GenerateUniqueEmail returns a unique email address for testing.
// Uses an atomic counter to ensure uniqueness across parallel test runs.
func GenerateUniqueEmail(prefix string) string {
	count := emailCounter.Add(1)
	return fmt.Sprintf("%s-%d@example.com", prefix, count)
}

// GenerateUniqueUsername returns a unique username for testing.
// Uses the same atomic counter as GenerateUniqueEmail to ensure uniqueness.
func GenerateUniqueUsername(prefix string) string {
	count := emailCounter.Add(1)
	return fmt.Sprintf("%s-%d", prefix, count)
}
