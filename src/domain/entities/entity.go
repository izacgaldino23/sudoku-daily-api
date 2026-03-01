package entities

import (
	"time"
)

var BoardSizes = map[int]int{4: 2, 6: 2, 9: 3}

type (
	BoardSize int

	Sudoku struct {
		ID    string
		Size  int
		Board [][]int
		Date  time.Time
	}
)
