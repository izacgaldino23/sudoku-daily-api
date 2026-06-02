package user

import (
	"context"
	"errors"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/repository"
)

type (
	RefreshTokenUseCase interface {
		Execute(ctx context.Context, tokenHash string) (accessToken, newRefreshToken string, err error)
	}

	userRefreshTokenUseCase struct {
		txManager        repository.TransactionManager
		refreshTokenRepo repository.RefreshTokenRepository
		tokenService     domain.TokenService
		userRepo         repository.UserRepository
	}
)

func NewUserRefreshTokenUseCase(
	txManager repository.TransactionManager,
	refreshTokenRepo repository.RefreshTokenRepository,
	tokenService domain.TokenService,
	userRepo repository.UserRepository,
) RefreshTokenUseCase {
	return &userRefreshTokenUseCase{
		txManager:        txManager,
		refreshTokenRepo: refreshTokenRepo,
		tokenService:     tokenService,
		userRepo:         userRepo,
	}
}

func (u *userRefreshTokenUseCase) Execute(ctx context.Context, tokenHash string) (string, string, error) {
	var newAccessToken, newRefreshTokenHash string

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		refreshToken, err := u.refreshTokenRepo.GetByToken(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, pkg.ErrRefreshTokenNotFound) {
				return pkg.ErrInvalidToken
			}
			return err
		}

		if refreshToken.Revoked {
			return pkg.ErrRefreshTokenRevoked
		}

		if refreshToken.ExpiresAt.Before(time.Now()) {
			return pkg.ErrRefreshTokenExpired
		}

		newAccessToken, err = u.tokenService.GenerateJWTToken(map[string]any{"user_id": refreshToken.UserID}, nil)
		if err != nil {
			return err
		}

		newTimezone := app_context.GetTimezoneFromContext(ctx)

		newRefreshToken, err := u.tokenService.GenerateRefreshToken(refreshToken.UserID, newTimezone)
		if err != nil {
			return err
		}

		// compare and update timezone if necessary
		if refreshToken.Timezone != newTimezone {
			if err := u.userRepo.UpdateTimezone(txCtx, refreshToken.UserID, newTimezone); err != nil {
				return err
			}
		}

		if err := u.refreshTokenRepo.Revoke(txCtx, refreshToken.UserID, tokenHash); err != nil {
			return err
		}

		if err := u.refreshTokenRepo.Create(txCtx, newRefreshToken); err != nil {
			return err
		}

		newRefreshTokenHash = newRefreshToken.Hash

		return nil
	}); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshTokenHash, nil
}
