package entities

import "time"

type Sudoku struct {
	ID    string
	Size  int
	Board [][]int
	Date  time.Time
}
