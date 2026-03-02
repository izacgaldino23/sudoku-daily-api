package helpers

import "slices"

func isLineValid(board [][]int, row, col, lines, cols int) bool {
	var (
		nonZero int
		size    = len(board)
	)

	line := make([]int, 0, lines*cols)
	for i := range lines {
		for j := range cols {
			line = append(line, board[row+i][col+j])
		}
	}

	// check total
	for _, v := range line {
		if v != 0 {
			nonZero++
		}
	}

	// check unique number from 1 to size
	if nonZero == size {
		for i := 1; i <= size; i++ {
			if !slices.Contains(line, i) {
				return false
			}
		}
	}

	// check repeat number
	for i := 0; i < len(line); i++ {
		for j := i + 1; j < len(line); j++ {
			if line[i] == 0 || line[j] == 0 {
				continue
			}

			if line[i] == line[j] {
				return false
			}
		}
	}

	return true
}
