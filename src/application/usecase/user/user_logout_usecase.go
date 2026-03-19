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
	refreshToken, err := u.refreshTokenRepo.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			return nil
		}
		return err
	}

	if refreshToken.Revoked {
		return nil
	} else if refreshToken.UserID != userID {
		return nil
	}

	err = u.refreshTokenRepo.Revoke(ctx, userID, refreshToken.Hash)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			return nil
		}
		return err
	}

	return nil
}
