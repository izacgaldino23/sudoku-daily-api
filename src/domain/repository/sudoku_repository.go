package repository

import (
	"context"
	"time"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
)

type SudokuRepository interface {
	Create(ctx context.Context, sudoku *entities.Sudoku) error
	GetByDateAndSize(ctx context.Context, date time.Time, size entities.BoardSize) (*entities.Sudoku, error)
	AddSolve(ctx context.Context, solve *entities.Solve) error
	GetTotalSolvedByUser(ctx context.Context, userID vo.UUID) (map[entities.BoardSize]int, error)
	GetTodaySolvedByUser(ctx context.Context, userID vo.UUID) ([]entities.Solve, error)
	GetBestTimesByUser(ctx context.Context, userID vo.UUID) ([]entities.Solve, error)
}
