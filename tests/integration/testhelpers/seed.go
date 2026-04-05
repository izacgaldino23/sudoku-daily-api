package testhelpers

import (
	"context"
	"fmt"
	"time"

	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/domain/vo"

	"github.com/uptrace/bun"
)

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
