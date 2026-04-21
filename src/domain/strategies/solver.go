package strategies

import (
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
)

type (
	solver struct {
		buffer []cell
	}

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
	empty := s.buffer[:0]

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

	return s.guess(board, empty, 0, 0)
}

func (s *solver) guess(board *entities.Board, empty []cell, left int, solutions int) int {
	if solutions >= 2 {
		return 2
	}	

	if left == 0 {
		return solutions + 1
	}

	best := -1
	bestCount := 999

	for i := 0; i < left; i++ {
		r := empty[i].row
		c := empty[i].col

		p := board.GetPossibleByPosition(r, c)
		cnt := p.Count()

		if cnt == 0 {
			return solutions
		}

		empty[i].possible = p

		if cnt < bestCount {
			best = i
			bestCount = cnt
			if cnt == 1 {
				break
			}
		}
	}

	empty[best], empty[left-1] = empty[left-1], empty[best]
	cur := empty[left-1]

	if cur.possible.Count() == 0 {
		return solutions
	}

	for _, n := range cur.possible.Values() {
		board.SetCell(cur.row, cur.col, n)

		solutions = s.guess(board, empty, left-1, solutions)
		board.SetCell(cur.row, cur.col, 0)

		if solutions > 1 {
			return 2
		}
	}

	return solutions
}
