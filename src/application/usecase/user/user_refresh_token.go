package user

import (
	"context"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
)

type (
	UserRefreshTokenUseCase interface {
		Execute(ctx context.Context, tokenHash string, userID vo.UUID) (accessToken string, err error)
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

func (u *userRefreshTokenUseCase) Execute(ctx context.Context, tokenHash string, userID vo.UUID) (string, error) {
	_, err := u.refreshTokenRepo.GetByToken(ctx, userID, tokenHash)
	if err != nil {
		return "", err
	}

	return u.tokenService.GenerateAccessToken(userID)
}
