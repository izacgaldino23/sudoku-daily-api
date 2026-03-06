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
		Execute(context.Context, *entities.User) (userData *entities.User, err error)
	}

	userLoginUseCase struct {
		txManager        repository.TransactionManager
		userRepo         repository.UserRepository
		refreshTokenRepo repository.RefreshTokenRepository
		passwordHasher   domain.PasswordHasher
		tokenService     domain.TokenService
	}
)

func NewUserLoginUseCase(
	txManager repository.TransactionManager,
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	passwordHasher domain.PasswordHasher,
	tokenService domain.TokenService,
) UserLoginUseCase {
	return &userLoginUseCase{
		txManager:        txManager,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordHasher:   passwordHasher,
		tokenService:     tokenService,
	}
}

func (u *userLoginUseCase) Execute(ctx context.Context, loginData *entities.User) (user *entities.User, err error) {
	err = u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		user, err = u.userRepo.GetByEmail(txCtx, loginData.Email.String())
		if err != nil {
			return err
		}

		if err := u.passwordHasher.Compare(*user.PasswordHash, *loginData.PasswordHash); err != nil {
			return pkg.ErrInvalidCredentials
		}

		user.Tokens = &entities.Tokens{}

		user.Tokens.AccessToken, err = u.tokenService.GenerateAccessToken(user.ID)
		if err != nil {
			return err
		}

		refreshToken, err := u.tokenService.GenerateRefreshToken(user.ID)
		if err != nil {
			return err
		}

		if err := u.refreshTokenRepo.Create(txCtx, refreshToken); err != nil {
			return err
		}

		user.Tokens.RefreshToken = refreshToken.Hash

		return nil
	})

	return
}
