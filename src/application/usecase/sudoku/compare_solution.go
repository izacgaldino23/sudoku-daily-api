package sudoku

import (
	"sudoku-daily-api/src/domain/entities"
)

func compareSolution(sudoku *entities.Sudoku, solution *entities.Solve) bool {
	board := sudoku.Solution.GetBoard()
	for i := 0; i < int(sudoku.Size); i++ {
		for j := 0; j < int(sudoku.Size); j++ {
			if board[i][j] != solution.Solution[i][j] {
				return false
			}
		}
	}

	return true
}