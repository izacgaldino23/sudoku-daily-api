package strategies

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"sudoku-daily-api/src/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestFillStrategy(t *testing.T) {
	f := NewFillStrategy()

	now := time.Now()
	fakeDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	r := rand.New(rand.NewSource(fakeDate.Unix()))

	for size := range entities.BoardSizes {
		t.Run(fmt.Sprintf("TestFillStrategy_%v", size), func(t *testing.T) {
			sudoku := entities.NewSudoku(size)
			sudoku.Date = fakeDate
			sudoku.Difficulty = entities.DifficultyMedium

			f.Fill(sudoku, r)

			emptyCells := 0
			for _, row := range sudoku.Board.GetBoard() {
				for _, cell := range row {
					if cell == 0 {
						emptyCells++
					}
				}
			}

			assert.Equal(t, 0, emptyCells)
		})
	}
}

func TestHideStrategy(t *testing.T) {
	h := NewHideStrategy()

	now := time.Now()
	fakeDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for size := range entities.BoardSizes {
		t.Run(fmt.Sprintf("TestHideStrategy_%v", size), func(t *testing.T) {
			r := rand.New(rand.NewSource(fakeDate.Unix()))
			sudoku := generateValidSudoku(size, r)
			sudoku.Date = fakeDate
			sudoku.Difficulty = entities.DifficultyMedium

			h.Hide(context.Background(), sudoku, r)

			emptyCells := 0
			for _, row := range sudoku.Board.GetBoard() {
				for _, cell := range row {
					if cell == 0 {
						emptyCells++
					}
				}
			}

			min, max := entities.GetClue(entities.BoardSize(size), entities.DifficultyMedium)

			var totalCells int = int(size * size)

			assert.GreaterOrEqual(t, totalCells-emptyCells, min, "min value for difficulty %v is %v", sudoku.Difficulty, min)
			assert.LessOrEqual(t, totalCells-emptyCells, max, "max value for difficulty %v is %v", sudoku.Difficulty, max)
		})
	}

}

func TestSolver(t *testing.T) {
	size := entities.BoardSize6

	board := [][]int{
		{0, 0, 3, 4, 5, 6},
		{4, 5, 6, 0, 0, 3},
		{3, 4, 5, 6, 0, 0},
		{6, 0, 0, 3, 4, 5},
		{0, 3, 4, 5, 6, 0},
		{5, 6, 0, 0, 3, 4},
	}

	sudoku := entities.NewSudoku(size)
	sudoku.Board = entities.NewFilledBoard(board)

	solver := newSolver()

	assert.Equal(t, 1, solver.Execute(&sudoku.Board))
}

func generateValidSudoku(size entities.BoardSize, r *rand.Rand) *entities.Sudoku {
	f := NewFillStrategy()

	sudoku := entities.NewSudoku(entities.BoardSize(size))

	f.Fill(sudoku, r)

	return sudoku
}

func TestGenerateComplete(t *testing.T) {
	hideStrategy := NewHideStrategy()
	fillStrategy := NewFillStrategy()

	now := time.Now()

	fakeDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	hardDate := time.Date(2026, 4, 18, 0, 0, 0, 0, time.UTC)

	dates := map[string]time.Time{
		"fake": fakeDate,
		"hard": hardDate,
	}

	for _, date := range dates {
		r := rand.New(rand.NewSource(date.Unix()))

		for size := range entities.BoardSizes {
			t.Run(fmt.Sprintf("TestGenerateComplete_%v_%v", date, size), func(t *testing.T) {

				sudoku := entities.NewSudoku(size)

				sudoku.Difficulty = entities.DifficultyMedium
				sudoku.Date = date

				filled := fillStrategy.Fill(sudoku, r)
				assert.True(t, filled)

				hidden := hideStrategy.Hide(context.Background(), sudoku, r)
				assert.True(t, hidden)

				filledCells := 0
				for _, row := range sudoku.Board.GetBoard() {
					for _, cell := range row {
						if cell != 0 {
							filledCells++
						}
					}
				}

				minClues, maxClues := entities.GetClue(sudoku.Size, sudoku.Difficulty)

				assert.GreaterOrEqual(t, filledCells, minClues, "min value for difficulty %v is %v", sudoku.Difficulty, minClues)
				assert.LessOrEqual(t, filledCells, maxClues, "max value for difficulty %v is %v", sudoku.Difficulty, maxClues)
			})
		}
	}

}
