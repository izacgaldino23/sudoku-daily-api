package user

import (
	"context"
	"errors"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
)

type (
	UserLogoutUseCase interface {
		Execute(ctx context.Context, userID vo.UUID, token string) error
	}

	userLogoutUseCase struct {
		refreshTokenRepo repository.RefreshTokenRepository
	}
)

func NewUserLogoutUseCase(refreshTokenRepo repository.RefreshTokenRepository) UserLogoutUseCase {
	return &userLogoutUseCase{refreshTokenRepo: refreshTokenRepo}
}

func (u *userLogoutUseCase) Execute(ctx context.Context, userID vo.UUID, token string) error {
	refreshToken, err := u.refreshTokenRepo.GetByToken(ctx, userID, token)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			return nil
		}
		return err
	}

	err = u.refreshTokenRepo.Revoke(ctx, userID, refreshToken.Hash)
	if err != nil {
		// TODO log the error in the database
		if errors.Is(err, pkg.ErrNotFound) {
			return nil
		}
		return err
	}

	return nil
}
