package helpers

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHideBacktracking(t *testing.T) {
	h := NewHideBacktracking()

	fakeDate, err := time.Parse("2006-01-02", "2022-01-01")
	assert.NoError(t, err)

	sudoku := entities.Sudoku{
		Board: [][]int{
			{1, 2, 3, 4},
			{3, 4, 1, 2},
			{2, 3, 4, 1},
			{4, 1, 2, 3},
		},
		Size:       4,
		Difficulty: entities.DifficultyMedium,
		Date:       fakeDate,
	}

	r := rand.New(rand.NewSource(sudoku.Date.Unix()))

	h.Hide(&sudoku, r)

	emptyCells := 0
	for _, row := range sudoku.Board {
		for _, cell := range row {
			if cell == 0 {
				emptyCells++
			}
		}
	}

	min, max := entities.GetClue(entities.BoardSize4, entities.DifficultyMedium)

	assert.GreaterOrEqual(t, emptyCells, min)
	assert.LessOrEqual(t, emptyCells, max)
}
