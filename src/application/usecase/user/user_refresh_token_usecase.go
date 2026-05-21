package user

import (
	"context"
	"errors"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
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
	}
)

func NewUserRefreshTokenUseCase(
	txManager repository.TransactionManager,
	refreshTokenRepo repository.RefreshTokenRepository,
	tokenService domain.TokenService,
) RefreshTokenUseCase {
	return &userRefreshTokenUseCase{
		txManager:        txManager,
		refreshTokenRepo: refreshTokenRepo,
		tokenService:     tokenService,
	}
}

func (u *userRefreshTokenUseCase) Execute(ctx context.Context, tokenHash string) (string, string, error) {
	refreshToken, err := u.refreshTokenRepo.GetByToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pkg.ErrRefreshTokenNotFound) {
			return "", "", pkg.ErrInvalidToken
		}
		return "", "", err
	}

	if refreshToken.Revoked {
		return "", "", pkg.ErrRefreshTokenRevoked
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return "", "", pkg.ErrRefreshTokenExpired
	}

	newAccessToken, err := u.tokenService.GenerateJWTToken(map[string]any{"user_id": refreshToken.UserID}, nil)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := u.tokenService.GenerateRefreshToken(refreshToken.UserID)
	if err != nil {
		return "", "", err
	}

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := u.refreshTokenRepo.Revoke(txCtx, refreshToken.UserID, tokenHash); err != nil {
			return err
		}

		if err := u.refreshTokenRepo.Create(txCtx, newRefreshToken); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken.Hash, nil
}
