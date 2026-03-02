package helpers

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/helpers"
)

type (
	hideBacktracking struct {
		baseNumber []int
	}
)

func NewHideBacktracking() helpers.HideBacktracking {
	return &hideBacktracking{}
}

func (h *hideBacktracking) Hide(board *entities.Sudoku, r *rand.Rand) {
	h.baseNumber = make([]int, board.Size)
	for i := 0; i < board.Size; i++ {
		h.baseNumber[i] = i + 1
	}

	h.hideNumbers(board, r)
}

func (s *hideBacktracking) hideNumbers(board *entities.Sudoku, r *rand.Rand) {
	difficulties := []entities.Difficulty{
		entities.DifficultyEasy,
		entities.DifficultyMedium,
		entities.DifficultyHard,
	}

	// get random difficulty
	difficulty := difficulties[r.Intn(len(difficulties))]
	min, max := entities.GetClue(entities.BoardSize(board.Size), difficulty)

	// get clue number between the range
	clueCount := rand.Intn(max-min+1) + min
	hideTotal := board.Size*board.Size - clueCount

	cellReference := make([][2]int, 0)
	for i := 0; i < board.Size; i++ {
		for j := 0; j < board.Size; j++ {
			cellReference = append(cellReference, [2]int{i, j})
		}
	}

	for {
		r.Shuffle(len(cellReference), func(i, j int) {
			cellReference[i], cellReference[j] = cellReference[j], cellReference[i]
		})

		// hide numbers
		if ok := s.hideCell(board, cellReference, 0, hideTotal); ok {
			break
		}
	}
}

func (s *hideBacktracking) hideCell(board *entities.Sudoku, toHide [][2]int, current, hideTotal int) bool {
	if current == hideTotal {
		return true
	}

	row, col := toHide[current][0], toHide[current][1]
	n := board.Board[row][col]

	board.Board[row][col] = 0

	// test solutions
	solutions := s.testSolutions(board)
	if solutions == 0 || solutions > 1 {
		board.Board[row][col] = n
		return false
	}

	return s.hideCell(board, toHide, current+1, hideTotal)
}

func (s *hideBacktracking) testSolutions(board *entities.Sudoku) (total int) {
	empty := make([][]int, 0)
	row := make([]uint8, board.Size)
	col := make([]uint8, board.Size)

	for i := 0; i < board.Size; i++ {
		for j := 0; j < board.Size; j++ {
			if board.Board[i][j] == 0 {
				empty = append(empty, []int{i, j})
			} else {
				row[i] = row[i] | (1 << board.Board[i][j])
				col[j] = col[j] | (1 << board.Board[i][j])
			}
		}
	}

	s.guess(board, row, col, empty, 0, &total)

	return
}

func (s *hideBacktracking) guess(board *entities.Sudoku, rowSum, colSum []uint8, empty [][]int, current int, solutions *int) {
	if current == len(empty) {
		*solutions++
		return
	}

	row, col := empty[current][0], empty[current][1]
	for i := 1; i <= board.Size; i++ {
		if (rowSum[row]&(1<<i)) == 0 && (colSum[col]&(1<<i)) == 0 {
			board.Board[row][col] = i
			rowSum[row] = rowSum[row] | (1 << i)
			colSum[col] = colSum[col] | (1 << i)

			s.guess(board, rowSum, colSum, empty, current+1, solutions)

			rowSum[row] = rowSum[row] ^ (1 << i)
			colSum[col] = colSum[col] ^ (1 << i)
			board.Board[row][col] = 0
		}
	}
}
