package services

import (
	"math/rand"
	"slices"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
)

type (
	grid struct {
		row      int
		col      int
		rowCount int
		colCount int
	}

	sudokuGenerator struct {
		baseNumber []int
	}
)

func (g *grid) isLastPosition(row, col int) bool {
	return row == g.row+g.rowCount-1 && col == g.col+g.colCount-1
}

func NewGenerator() domain.Generator {
	return &sudokuGenerator{}
}

func (s *sudokuGenerator) GenerateDaily(size int, difficulty string, seed int64) *entities.Sudoku {
	sum := 0
	for i := 0; i < size; i++ {
		sum += i + 1
	}

	board := &entities.Sudoku{
		Size:  size,
		Board: make([][]int, size),
	}

	s.baseNumber = make([]int, board.Size)
	for i := 0; i < board.Size; i++ {
		s.baseNumber[i] = i + 1
	}

	// initialize board
	for i := range board.Board {
		board.Board[i] = make([]int, size)
	}

	s.generateTiles(board)

	return board
}

func (s *sudokuGenerator) generateTiles(board *entities.Sudoku) {
	r := rand.New(rand.NewSource(board.Date.Unix()))

	var (
		chosen []int
	)

	grids := s.getGrids(board.Size)

	s.iterateCell(board, 0, 0, chosen, r, grids)
}

func (s *sudokuGenerator) iterateCell(board *entities.Sudoku, currentRow, currentCol int, chosen []int, r *rand.Rand, grids []grid) bool {
	var currentDecision []int

	if len(chosen) == 0 {
		currentDecision = make([]int, board.Size-len(chosen))
		copy(currentDecision, s.baseNumber)
	} else {
		for i := range board.Size {
			if !slices.Contains(chosen, i+1) {
				currentDecision = append(currentDecision, i+1)
			}
		}
	}

	// shuffle numbers
	if len(currentDecision) > 1 {
		r.Shuffle(len(currentDecision), func(i, j int) {
			currentDecision[i], currentDecision[j] = currentDecision[j], currentDecision[i]
		})
	}

	for i := range currentDecision {
		n := currentDecision[i]
		board.Board[currentRow][currentCol] = n

		// validate line
		if !s.isLineValid(board.Board, currentRow, 0, 1, board.Size) {
			continue
		}

		// validate columns
		if !s.isLineValid(board.Board, 0, currentCol, board.Size, 1) {
			continue
		}

		// validate grid
		valid := true
		for _, grid := range grids {
			if grid.isLastPosition(currentRow, currentCol) {
				if !s.isLineValid(board.Board, grid.row, grid.col, grid.rowCount, grid.colCount) {
					valid = false
				}
			}
		}
		if !valid {
			continue
		}

		if currentCol == board.Size-1 && currentRow == board.Size-1 {
			return true
		}

		if currentCol == board.Size-1 {
			if s.iterateCell(board, currentRow+1, 0, []int{}, r, grids) {
				return true
			}
		} else {
			// call the next cell
			if s.iterateCell(board, currentRow, currentCol+1, append(chosen, n), r, grids) {
				return true
			}
		}
	}

	board.Board[currentRow][currentCol] = 0

	return false
}

// func (s *sudokuGenerator) isLineValid(line []int, total, size int) bool {
func (s *sudokuGenerator) isLineValid(board [][]int, row, col, lines, cols int) bool {
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

func (s *sudokuGenerator) getGrids(size int) []grid {
	grids := make([]grid, 0)

	gridRows := entities.BoardSizes[size]
	gridCols := size / gridRows

	for i := 0; i < gridRows; i++ {
		for j := 0; j < gridCols; j++ {
			grids = append(grids, grid{
				row:      i * gridRows,
				col:      j * gridCols,
				rowCount: gridRows,
				colCount: gridCols,
			})
		}
	}

	return grids
}
