package helpers

import (
	"sudoku-daily-api/src/domain/vo"
)

func isLineValid(board [][]int, row, col, lines, cols int) bool {
	var (
		nonZero int
		size    = len(board)
	)

	// line := make([]int, 0, lines*cols)
	var line vo.Binary

	for i := range lines {
		for j := range cols {
			v := board[row+i][col+j]

			if v != 0 && line.Contains(v) {
				return false
			}

			line.Add(v)
			if v != 0 {
				nonZero++
			}
		}
	}

	// check unique number from 1 to size
	if nonZero == size {
		for i := 1; i <= size; i++ {
			if !line.Contains(i) {
				return false
			}
		}
	}

	return true
}
