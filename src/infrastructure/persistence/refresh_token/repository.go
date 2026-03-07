package refresh_token

import (
	"context"
	"database/sql"
	"errors"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/entities"
	repository "sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/src/infrastructure/persistence/tx"

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
		return err
	}

	if a, err := res.RowsAffected(); err != nil {
		return err
	} else if a == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, userID vo.UUID, token string) (*entities.RefreshToken, error) {
	var refreshTokenModel RefreshToken

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Model(&refreshTokenModel).
		Where("token_hash = ? AND user_id = ?", token, userID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}
		return nil, err
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

	return err
}
