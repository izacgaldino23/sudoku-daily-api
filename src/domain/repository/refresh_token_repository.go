package repository

import (
	"context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	GetByToken(ctx context.Context, userID vo.UUID, token string) (*entities.RefreshToken, error)
	Revoke(ctx context.Context, userID vo.UUID, token string) error
}
