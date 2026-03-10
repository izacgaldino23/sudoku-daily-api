package sudoku

import (
	"math"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/uptrace/bun"
)

type (
	Sudoku struct {
		bun.BaseModel `bun:"table:sudokus"`

		ID         string    `bun:"id,pk"`
		Size       int       `bun:",notnull"`
		Difficulty string    `bun:",notnull"`
		Board      []byte    `bun:"type:,notnull"`
		Solution   []byte    `bun:"type:,notnull"`
		Date       time.Time `bun:"type:date,notnull"`
	}

	Solve struct {
		bun.BaseModel `bun:"table:solves"`

		ID          string    `bun:"id,pk"`
		SudokuID    string    `bun:"sudoku_id,notnull"`
		UserID      string    `bun:"user_id,notnull"`
		StartedAt   time.Time `bun:"type:timestamp,notnull"`
		CompletedAt time.Time `bun:"type:timestamp,notnull"`
		CreatedAt   time.Time `bun:"type:timestamp,notnull"`
	}
)

func (s *Sudoku) FromDomain(sudoku *entities.Sudoku) {
	s.ID = sudoku.ID.String()
	s.Size = sudoku.GetSize()
	s.Difficulty = string(sudoku.Difficulty)
	s.Board = boardFromDomain(&sudoku.Board)
	s.Solution = boardFromDomain(&sudoku.Solution)
	s.Date = sudoku.Date
}

func (s *Sudoku) ToDomain() *entities.Sudoku {
	return &entities.Sudoku{
		ID:         vo.UUID(s.ID),
		Size:       entities.BoardSize(s.Size),
		Difficulty: entities.Difficulty(s.Difficulty),
		Board:      boardToDomain(s.Board),
		Solution:   boardToDomain(s.Solution),
		Date:       s.Date,
	}
}

func (s *Solve) FromDomain(solve *entities.Solve) {
	s.ID = string(solve.ID)
	s.SudokuID = string(solve.SudokuID)
	s.UserID = string(solve.UserID)
	s.StartedAt = solve.StartedAt
	s.CompletedAt = solve.CompletedAt
	s.CreatedAt = solve.CreatedAt
}

func boardToDomain(boardData []byte) entities.Board {
	size := int(math.Sqrt(float64(len(boardData))))

	boardFilled := make([][]int, size)

	for i := 0; i < size; i++ {
		row := make([]int, size)
		for j := 0; j < size; j++ {
			row[j] = int(boardData[i*size+j])
		}
		boardFilled[i] = row
	}

	return entities.NewFilledBoard(boardFilled)
}

func boardFromDomain(board *entities.Board) []byte {
	size := board.GetSize()
	linearBoard := make([]byte, 0, size*size)

	for _, row := range board.GetBoard() {
		for _, val := range row {
			linearBoard = append(linearBoard, byte(val))
		}
	}

	return linearBoard
}
