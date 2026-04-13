package refresh_token

import (
	"context"
	"database/sql"
	"errors"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/entities"
	repository "sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/src/infrastructure/persistence/database/tx"

	"github.com/uptrace/bun"
)

type (
	refreshTokenRepository struct {
		txManager *tx.Manager
		db        *bun.DB
	}
)

func NewRepository(db *bun.DB) repository.RefreshTokenRepository {
	return &refreshTokenRepository{
		db:        db,
		txManager: tx.NewManager(db),
	}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	refreshToken := NewModel(token)

	res, err := r.txManager.GetExecutor(ctx).
		NewInsert().
		Model(refreshToken).
		Exec(ctx)
	if err != nil {
		return r.txManager.HandleError(ctx, err)
	}

	if a, err := res.RowsAffected(); err != nil {
		return r.txManager.HandleError(ctx, err)
	} else if a == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	var refreshTokenModel RefreshToken

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Model(&refreshTokenModel).
		Where("token_hash = ?", token).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrRefreshTokenNotFound
		}
		return nil, r.txManager.HandleError(ctx, err)
	}

	return refreshTokenModel.ToDomain(), nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, userID vo.UUID, token string) error {
	_, err := r.txManager.GetExecutor(ctx).
		NewUpdate().
		Model(&RefreshToken{}).
		Where("token_hash = ? AND user_id = ?", token, userID).
		Set("revoked = ?", true).
		Exec(ctx)

	return r.txManager.HandleError(ctx, err)
}
