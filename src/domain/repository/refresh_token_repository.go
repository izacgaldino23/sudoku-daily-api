package repository

import (
	"context"
	"sudoku-daily-api/src/domain/entities"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
}
