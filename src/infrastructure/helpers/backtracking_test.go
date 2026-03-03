package helpers

import (
	"fmt"
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFillBacktracking(t *testing.T) {
	f := NewFillBacktracking()

	fakeDate, err := time.Parse("2006-01-02", "2022-01-01")
	assert.NoError(t, err)

	r := rand.New(rand.NewSource(fakeDate.Unix()))

	for size := range entities.BoardSizes {
		t.Run(fmt.Sprintf("TestFillBacktracking_%v", size), func(t *testing.T) {
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

func TestHideBacktracking(t *testing.T) {
	h := NewHideBacktracking()

	fakeDate, err := time.Parse("2006-01-02", "2022-01-01")
	assert.NoError(t, err)

	size := 6

	r := rand.New(rand.NewSource(fakeDate.Unix()))
	sudoku := generateValidSudoku(size, r)
	sudoku.Date = fakeDate
	sudoku.Difficulty = entities.DifficultyMedium

	h.Hide(sudoku, r)

	emptyCells := 0
	for _, row := range sudoku.Board.GetBoard() {
		for _, cell := range row {
			if cell == 0 {
				emptyCells++
			}
		}
	}

	min, max := entities.GetClue(entities.BoardSize(size), entities.DifficultyMedium)

	totalCells := size * size

	assert.GreaterOrEqual(t, totalCells-emptyCells, min)
	assert.LessOrEqual(t, totalCells-emptyCells, max)
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

	solver := NewSolver()

	assert.Equal(t, 2, solver.Execute(sudoku))
}

func generateValidSudoku(size int, r *rand.Rand) *entities.Sudoku {
	f := NewFillBacktracking()

	sudoku := entities.NewSudoku(entities.BoardSize(size))

	f.Fill(sudoku, r)

	return sudoku
}

func TestGenerateComplete(t *testing.T) {
	hideBacktracking := NewHideBacktracking()
	fillBacktracking := NewFillBacktracking()

	size := entities.BoardSize4

	fakeDate, err := time.Parse("2006-01-02", "2022-01-01")
	assert.NoError(t, err)

	r := rand.New(rand.NewSource(fakeDate.Unix()))

	sudoku := entities.NewSudoku(size)

	sudoku.Difficulty = entities.DifficultyMedium
	sudoku.Date = fakeDate

	fillBacktracking.Fill(sudoku, r)
	hideBacktracking.Hide(sudoku, r)

	emptyCells := 0
	for _, row := range sudoku.Board.GetBoard() {
		for _, cell := range row {
			if cell == 0 {
				emptyCells++
			}
		}
	}

	min, max := entities.GetClue(sudoku.Size, sudoku.Difficulty)

	totalCells := int(size * size)

	assert.GreaterOrEqual(t, totalCells-emptyCells, min, "min value for difficulty %v is %v", sudoku.Difficulty, min)
	assert.LessOrEqual(t, totalCells-emptyCells, max, "max value for difficulty %v is %v", sudoku.Difficulty, max)
}
