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
	UserRefreshTokenUseCase interface {
		Execute(ctx context.Context, tokenHash string) (accessToken string, err error)
	}

	userRefreshTokenUseCase struct {
		refreshTokenRepo repository.RefreshTokenRepository
		tokenService     domain.TokenService
	}
)

func NewUserRefreshTokenUseCase(
	refreshTokenRepo repository.RefreshTokenRepository,
	tokenService domain.TokenService,
) UserRefreshTokenUseCase {
	return &userRefreshTokenUseCase{
		refreshTokenRepo: refreshTokenRepo,
		tokenService:     tokenService,
	}
}

func (u *userRefreshTokenUseCase) Execute(ctx context.Context, tokenHash string) (string, error) {
	refreshToken, err := u.refreshTokenRepo.GetByToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			return "", pkg.ErrInvalidToken
		}
		return "", err
	}

	if refreshToken.Revoked {
		return "", pkg.ErrRefreshTokenRevoked
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return "", pkg.ErrRefreshTokenExpired
	}

	return u.tokenService.GenerateJWTToken(map[string]any{"user_id": refreshToken.UserID}, nil)
}
