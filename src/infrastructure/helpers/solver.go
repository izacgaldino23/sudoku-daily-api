package helpers

import (
	"slices"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
)

type (
	Solver struct {}

	cell struct {
		row      int
		col      int
		possible vo.Binary
	}
)

func NewSolver() *Solver {
	return &Solver{}
}

func (s *Solver) Execute(board *entities.Sudoku) int {
	empty := make([]cell, 0)

	for i := 0; i < board.GetSize(); i++ {
		for j := 0; j < board.GetSize(); j++ {
			if board.Board.GetCell(i, j) == 0 {
				empty = append(empty, cell{
					row:      i,
					col:      j,
				})
			}
		}
	}

	for i := range empty {
		v := &empty[i]
		v.possible = vo.Binary(board.Board.GetPossibleByPosition(v.row, v.col))
	}

	// put the item with less possible values first
	slices.SortFunc(empty, func(a, b cell) int {
		return a.possible.Count() - b.possible.Count()
	})

	return s.guess(board, empty, 0, 0)
}

func (s *Solver) guess(board *entities.Sudoku, empty []cell, current int, solutions int) int {
	if current == len(empty) {
		return solutions + 1
	}

	row, col := empty[current].row, empty[current].col

	for _, n := range empty[current].possible.Values() {
		if !board.Board.HasNumber(row, col, n) {
			board.Board.SetCell(row, col, n)

			if v := s.guess(board, empty, current+1, solutions); v > 1 {
				return v
			} else {
				solutions = v
			}

			board.Board.SetCell(row, col, 0)
		}
	}

	return solutions
}
