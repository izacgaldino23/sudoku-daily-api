package repository

import (
	"context"
	"sudoku-daily-api/src/domain/entities"
	"time"
)

type (
	SudokuRepository interface {
		Create(ctx context.Context, sudoku *entities.Sudoku) error
		GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error)
	}

	UserRepository interface {
		Create(ctx context.Context, user *entities.User) error
		GetByEmail(ctx context.Context, email string) (*entities.User, error)
	}

	RefreshTokenRepository interface {
		Create(ctx context.Context, token *entities.RefreshToken) error
	}
)
