package strategies

import (
	"slices"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
)

type (
	solver struct{}

	cell struct {
		row      int
		col      int
		possible vo.Binary
	}
)

func newSolver() *solver {
	return &solver{}
}

func (s *solver) Execute(board *entities.Board) int {
	empty := make([]cell, 0)

	full := board.GetFullCount()

	for i := 0; i < board.GetSize(); i++ {
		if board.RowCount[i] == full {
			continue
		}

		for j := 0; j < board.GetSize(); j++ {
			if board.GetCell(i, j) == 0 {
				empty = append(empty, cell{
					row:      i,
					col:      j,
					possible: board.GetPossibleByPosition(i, j),
				})
			}
		}
	}

	slices.SortFunc(empty, func(a, b cell) int {
		return a.possible.Count() - b.possible.Count()
	})

	return s.guess(board, empty, 0, 0)
}

func (s *solver) guess(board *entities.Board, empty []cell, current int, solutions int) int {
	if current == len(empty) {
		return solutions + 1
	}

	row, col := empty[current].row, empty[current].col
	possibilities := board.GetPossibleByPosition(row, col)

	if possibilities.Count() == 0 {
		return solutions
	}

	for _, n := range possibilities.Values() {
		if !board.HasNumber(row, col, n) {
			board.SetCell(row, col, n)

			v := s.guess(board, empty, current+1, solutions)
			board.SetCell(row, col, 0)
			if v > 1 {
				return v
			} else {
				solutions = v
			}

			board.SetCell(row, col, 0)
		}
	}

	return solutions
}
