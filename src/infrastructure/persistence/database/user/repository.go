package user

import (
	"context"
	"database/sql"
	"errors"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/entities"
	repository "sudoku-daily-api/src/domain/repository"
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
		return err
	}

	_, err = result.RowsAffected()
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var userResp User

	err := r.txManager.GetExecutor(ctx).NewSelect().Model(&userResp).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}
		return nil, err
	}

	return userResp.ToDomain(), nil
}
