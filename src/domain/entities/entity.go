package entities

import (
	"time"
)

const (
	BoardSize4 BoardSize = 4
	BoardSize6 BoardSize = 6
	BoardSize9 BoardSize = 9

	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

var (
	// the map value represents the number of horizontal grids or lines in each grid
	BoardSizes = map[BoardSize]int{
		BoardSize4: 2,
		BoardSize6: 2,
		BoardSize9: 3,
	}

	minimumClues = map[BoardSize]map[Difficulty]Clue{
		BoardSize4: {
			DifficultyEasy:   Clue{minimum: 8, maximum: 10},
			DifficultyMedium: Clue{minimum: 6, maximum: 7},
			DifficultyHard:   Clue{minimum: 4, maximum: 5},
		},
		BoardSize6: {
			DifficultyEasy:   Clue{minimum: 18, maximum: 22},
			DifficultyMedium: Clue{minimum: 14, maximum: 17},
			DifficultyHard:   Clue{minimum: 10, maximum: 13},
		},
		BoardSize9: {
			DifficultyEasy:   Clue{minimum: 36, maximum: 45},
			DifficultyMedium: Clue{minimum: 30, maximum: 35},
			DifficultyHard:   Clue{minimum: 22, maximum: 28},
		},
	}
)

type (
	BoardSize int

	Difficulty string

	Clue struct {
		minimum int
		maximum int
	}

	Sudoku struct {
		ID         string
		Size       int
		Board      [][]int
		Difficulty Difficulty
		Date       time.Time
	}

	Grid struct {
		Row      int
		Col      int
		RowCount int
		ColCount int
	}
)

func (s *Sudoku) GetGrids(size int) []Grid {
	grids := make([]Grid, 0)

	gridRows := BoardSizes[BoardSize(size)]
	gridCols := size / gridRows

	for i := 0; i < gridRows; i++ {
		for j := 0; j < gridCols; j++ {
			grids = append(grids, Grid{
				Row:      i * gridRows,
				Col:      j * gridCols,
				RowCount: gridRows,
				ColCount: gridCols,
			})
		}
	}

	return grids
}

func (g *Grid) IsLastPosition(row, col int) bool {
	return row == g.Row+g.RowCount-1 && col == g.Col+g.ColCount-1
}

func GetClue(boardSize BoardSize, difficulty Difficulty) (int, int) {
	clue := minimumClues[boardSize][difficulty]

	return clue.minimum, clue.maximum
}
