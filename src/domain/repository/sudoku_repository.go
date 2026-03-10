package repository

import (
	"context"
	"sudoku-daily-api/src/domain/entities"
	"time"
)

type SudokuRepository interface {
	Create(ctx context.Context, sudoku *entities.Sudoku) error
	GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error)
	AddSolve(ctx context.Context, solve *entities.Solve) error
}
