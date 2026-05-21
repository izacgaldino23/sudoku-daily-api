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
	AddAttempt(ctx context.Context, attempt *entities.Solve) error
	MarkAsSolved(ctx context.Context, solve *entities.Solve, completedAt time.Time) error
	GetTotalSolvedByUser(ctx context.Context, userID vo.UUID) (map[entities.BoardSize]int, error)
	GetTodaySolvedByUser(ctx context.Context, userID vo.UUID) ([]entities.Solve, error)
	GetBestTimesByUser(ctx context.Context, userID vo.UUID) ([]entities.Solve, error)
	GetDailyLeaderboard(ctx context.Context, sudokuID vo.UUID, limit, offset int) ([]entities.Solve, bool, error)
	GetAllTimeBestLeaderboard(ctx context.Context, size entities.BoardSize, limit, offset int) ([]entities.Solve, bool, error)
	GetSolveByIDAndUser(ctx context.Context, userID vo.UUID, sudokuID vo.UUID) (*entities.Solve, error)
	GetSolvesByUserAndDate(ctx context.Context, userID vo.UUID, date time.Time) ([]entities.Solve, error)
	RemoveUnfinishedAttempts(ctx context.Context, date time.Time) (int64, error)
}
