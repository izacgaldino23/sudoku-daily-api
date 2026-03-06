package persistence

import (
	"context"
	"database/sql"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"

	"github.com/uptrace/bun"
)

type (
	refreshTokenRepository struct {
		transactionManager
		db *bun.DB
	}
)

func NewRefreshTokenRepository(db *bun.DB) repository.RefreshTokenRepository {
	return &refreshTokenRepository{
		db:                 db,
		transactionManager: transactionManager{db: db},
	}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	refreshToken := &RefreshToken{
		ID:        vo.NewUUID().String(),
		UserID:    token.UserID.String(),
		TokenHash: token.Hash,
		ExpiresAt: token.ExpiresAt,
		Revoked:   false,
	}

	res, err := r.GetExecutor(ctx).
		NewInsert().
		Model(refreshToken).
		Exec(ctx)
	if err != nil {
		return err
	}

	if a, err := res.RowsAffected(); err != nil {
		return err
	} else if a == 0 {
		return sql.ErrNoRows
	}

	return nil
}
