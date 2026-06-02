package user

import (
	"context"
	"database/sql"
	"errors"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/src/infrastructure/persistence/database/tx"

	"github.com/uptrace/bun"
)

type (
	userRepository struct {
		txManager *tx.Manager
		db        *bun.DB
	}
)

func NewRepository(db *bun.DB) repository.UserRepository {
	return &userRepository{
		db:        db,
		txManager: tx.NewManager(db),
	}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	var userModel = &User{}
	userModel.FromDomain(user)

	result, err := r.txManager.GetExecutor(ctx).
		NewInsert().
		Model(userModel).
		Exec(ctx)
	if err != nil {
		return r.txManager.HandleError(ctx, err)
	}

	_, err = result.RowsAffected()
	return r.txManager.HandleError(ctx, err)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var userResp User

	err := r.txManager.GetExecutor(ctx).NewSelect().Model(&userResp).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrUserNotFound
		}
		return nil, r.txManager.HandleError(ctx, err)
	}

	return userResp.ToDomain(), nil
}

func (r *userRepository) UpdateTimezone(ctx context.Context, userID vo.UUID, timezone string) error {
	var userModel = &User{}

	_, err := r.txManager.GetExecutor(ctx).
		NewUpdate().
		Model(userModel).
		Set("timezone = ?", timezone).
		Where("id = ?", userID).
		Exec(ctx)

	return r.txManager.HandleError(ctx, err)
}
