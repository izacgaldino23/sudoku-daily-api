package entities

import (
	"sudoku-daily-api/src/domain/vo"
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
			DifficultyEasy:   Clue{minimum: 9, maximum: 10},
			DifficultyMedium: Clue{minimum: 7, maximum: 8},
			DifficultyHard:   Clue{minimum: 5, maximum: 6},
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
		Size       BoardSize
		Board      Board
		Difficulty Difficulty
		Date       time.Time
	}

	Board struct {
		cells     [][]int
		RowCount  []vo.Binary
		ColCount  []vo.Binary
		GridCount []vo.Binary
		fullCount vo.Binary
	}

	Grid struct {
		Row      int
		Col      int
		RowCount int
		ColCount int
	}
)

func NewSudoku(size BoardSize) *Sudoku {
	return &Sudoku{Size: size, Board: newBoard(int(size))}
}

func newBoard(size int) Board {
	cells := make([][]int, size)

	for i := range cells {
		cells[i] = make([]int, size)
	}

	return Board{
		cells:     cells,
		RowCount:  make([]vo.Binary, size),
		ColCount:  make([]vo.Binary, size),
		GridCount: make([]vo.Binary, size*size),
		fullCount: vo.NewFullBinary(size),
	}
}

func NewFilledBoard(values [][]int) Board {
	board := newBoard(len(values))

	for i := range values {
		for j := range values[i] {
			board.SetCell(i, j, values[i][j])
		}
	}

	return board
}

func (s *Sudoku) GetSize() int {
	return int(s.Size)
}

func (b *Board) SetCell(row, col, value int) {

	if value == 0 {
		n := b.cells[row][col]

		b.RowCount[row].Remove(n)
		b.ColCount[col].Remove(n)
		b.GridCount[b.GetGridByPosition(row, col)].Remove(n)
	} else {
		b.RowCount[row].Add(value)
		b.ColCount[col].Add(value)
		b.GridCount[b.GetGridByPosition(row, col)].Add(value)
	}

	b.cells[row][col] = value
}

func (b *Board) GetCell(row, col int) int {
	return b.cells[row][col]
}

func (b *Board) GetSize() int {
	return len(b.cells)
}

func (b *Board) GetBoard() [][]int {
	return b.cells
}

func (b *Board) GetPossibleByPosition(row, col int) vo.Binary {
	var possible = b.GetFullCount()
	var current vo.Binary

	current.Union(b.RowCount[row], b.ColCount[col], b.GridCount[b.GetGridByPosition(row, col)])

	return current.Missing(possible)
	// possible.Sub(b.RowCount[row])
	// possible.Sub(b.ColCount[col])
	// possible.Sub(b.GridCount[b.GetGridByPosition(row, col)])

	// return possible
}

func (b *Board) GetRowMissingNumbers(row int) []int {
	missing := b.RowCount[row].Missing(b.GetFullCount())

	return missing.Values()
}

func (b *Board) GetColMissingNumbers(col int) []int {
	missing := b.ColCount[col].Missing(b.GetFullCount())

	return missing.Values()
}

func (b *Board) GetGridMissingNumbers(row, col int) []int {
	missing := b.GridCount[b.GetGridByPosition(row, col)].Missing(b.GetFullCount())

	return missing.Values()
}

func (b *Board) GetGridByPosition(currentRow, currentCol int) int {
	size := b.GetSize()
	rowsPerGrid := BoardSizes[BoardSize(size)]
	colsPerGrid := size / rowsPerGrid

	numGridsWide := size / colsPerGrid

	gridX := currentCol / colsPerGrid
	gridY := currentRow / rowsPerGrid

	return gridX + (gridY * numGridsWide)
}

func (b *Board) HasNumber(row, col, number int) bool {
	return b.RowCount[row].Contains(number) || b.ColCount[col].Contains(number) || b.GridCount[b.GetGridByPosition(row, col)].Contains(number)
}

func (b *Board) GetFullCount() vo.Binary {
	if b.fullCount == 0 {
		b.fullCount = vo.NewFullBinary(b.GetSize())
	}
	return b.fullCount
}

func GetClue(boardSize BoardSize, difficulty Difficulty) (int, int) {
	clue := minimumClues[boardSize][difficulty]

	return clue.minimum, clue.maximum
}
