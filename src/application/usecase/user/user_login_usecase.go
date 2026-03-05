package user

import (
	"context"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
)

type (
	UserLoginUseCase interface {
		Execute(context.Context, *entities.User) (userData *entities.User, accessToken string, refreshToken string, err error)
	}

	userLoginUseCase struct {
		userRepo       repository.UserRepository
		passwordHasher domain.PasswordHasher
		tokenService   domain.TokenService
	}
)

func NewUserLoginUseCase(
	userRepo repository.UserRepository,
	passwordHasher domain.PasswordHasher,
	tokenService domain.TokenService,
) UserLoginUseCase {
	return &userLoginUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

func (u *userLoginUseCase) Execute(ctx context.Context, loginData *entities.User) (*entities.User, string, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, loginData.Email.String())
	if err != nil {
		return nil, "", "", err
	}

	if err := u.passwordHasher.Compare(*user.PasswordHash, *loginData.PasswordHash); err != nil {
		return nil, "", "", pkg.ErrInvalidCredentials
	}

	accessToken, err := u.tokenService.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := u.tokenService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
