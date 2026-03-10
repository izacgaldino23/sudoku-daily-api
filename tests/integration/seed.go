package integration

import (
	"context"
	"math"
	"sudoku-daily-api/pkg/database"
	"time"

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
			ID:         "00000000-0000-0000-0000-000000000001",
			Size:       9,
			Difficulty: "easy",
			Board:      []byte{0, 0, 0, 0, 9, 4, 0, 3, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 1, 8, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 2, 3, 0, 0, 0, 0, 0, 5, 0, 0, 0, 8, 0, 7, 0, 0, 4, 0, 0, 0, 0, 0, 0, 3, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 7, 4, 0, 0, 0, 0, 0, 0, 0},
			Solution:   []byte{5, 8, 7, 6, 9, 4, 1, 3, 2, 6, 1, 4, 3, 2, 7, 9, 5, 8, 9, 2, 3, 8, 1, 5, 4, 7, 6, 3, 7, 8, 2, 4, 9, 6, 1, 5, 4, 5, 9, 7, 3, 6, 8, 2, 1, 1, 4, 6, 5, 7, 3, 2, 8, 9, 2, 3, 1, 9, 6, 8, 5, 7, 4, 7, 9, 5, 1, 3, 6, 8, 2, 4, 8, 6, 2, 4, 5, 1, 3, 9, 7, 1, 5, 3, 7, 2, 9, 4, 6, 8},
			Date:       today,
		},
		{
			ID:         "00000000-0000-0000-0000-000000000002",
			Size:       4,
			Difficulty: "easy",
			Board:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Solution:   []byte{1, 2, 3, 4, 3, 4, 1, 2, 2, 1, 4, 3, 4, 3, 2, 1},
			Date:       today,
		},
		{
			ID:         "00000000-0000-0000-0000-000000000003",
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

func GetSudokuSolution(size int) ([][]int, error) {
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
