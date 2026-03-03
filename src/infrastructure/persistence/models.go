package persistence

import (
	"sudoku-daily-api/src/domain/entities"
	"time"

	"github.com/uptrace/bun"
)

type (
	Sudoku struct {
		bun.BaseModel `bun:"table:sudoku"`

		ID    string    `bun:"id,pk"`
		Size  int       `bun:",notnull"`
		Board [][]int   `bun:"type:jsonb,notnull"`
		Date  time.Time `bun:",notnull"`
	}
)

func (s *Sudoku) ToDomain() *entities.Sudoku {
	return &entities.Sudoku{
		ID:    s.ID,
		Size:  entities.BoardSize(s.Size),
		Board: entities.NewFilledBoard(s.Board),
		Date:  s.Date,
	}
}

func (s *Sudoku) FromDomain(sudoku *entities.Sudoku) {
	s.ID = sudoku.ID
	s.Size = sudoku.GetSize()
	s.Board = sudoku.Board.GetBoard()
	s.Date = sudoku.Date
}
