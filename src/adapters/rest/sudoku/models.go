package sudoku

import (
	"sudoku-daily-api/src/domain/entities"
	"time"
)

type (
	GetDailySudokuResponse struct {
		ID    string    `json:"id"`
		Size  int       `json:"size"`
		Board [][]int   `json:"board"`
		Date  time.Time `json:"date"`
	}
)

func (g *GetDailySudokuResponse) FromDomain(s *entities.Sudoku) {
	g.ID = s.ID
	g.Size = s.Size
	g.Board = s.Board
	g.Date = s.Date
}