package sudoku

import (
	"sudoku-daily-api/src/domain/entities"
	"time"
)

type (
	GetDailySudokuResponse struct {
		ID    string `json:"id"`
		Size  int    `json:"size"`
		Board []Cell `json:"board"`
		Date  string `json:"date"`
	}

	Cell struct {
		Row   int `json:"row"`
		Col   int `json:"col"`
		Value int `json:"value"`
	}
)

func (g *GetDailySudokuResponse) FromDomain(s *entities.Sudoku) {
	g.ID = s.ID
	g.Size = s.GetSize()
	g.Board = BoardFromDomain(s.Board)
	g.Date = s.Date.Format(time.DateOnly)
}

func BoardFromDomain(board entities.Board) []Cell {
	var cells = make([]Cell, 0)

	for i, row := range board.GetBoard() {
		for j := range row {
			val := board.GetCell(i, j)
			if val == 0 {
				continue
			}
			cells = append(cells, Cell{
				Row:   i,
				Col:   j,
				Value: val,
			})
		}
	}

	return cells
}
