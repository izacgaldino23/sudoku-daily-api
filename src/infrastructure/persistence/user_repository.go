package persistence

import (
	"context"
	"database/sql"
	"errors"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"

	"github.com/uptrace/bun"
)

type (
	userRepository struct {
		transactionManager
		db *bun.DB
	}
)

func NewUserRepository(db *bun.DB) repository.UserRepository {
	return &userRepository{
		db:                 db,
		transactionManager: transactionManager{db: db},
	}
}

func (u *userRepository) Create(ctx context.Context, user *entities.User) error {
	var userModel = &User{}
	userModel.FromDomain(user)

	result, err := u.GetExecutor(ctx).
		NewInsert().
		Model(userModel).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var userResp User

	err := u.GetExecutor(ctx).NewSelect().Model(&userResp).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}
		return nil, err
	}

	return userResp.ToDomain(), nil
}
