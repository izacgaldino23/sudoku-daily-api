package user

import (
	"context"
	"errors"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
)

type (
	LoginUseCase interface {
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
) LoginUseCase {
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
			if errors.Is(err, pkg.ErrUserNotFound) {
				return pkg.ErrInvalidCredentials
			}
			return err
		}

		if err := u.passwordHasher.Compare(loginData.PasswordHash, user.PasswordHash); err != nil {
			return pkg.ErrInvalidCredentials
		}

		// compare and update timezone if necessary
		newTimezone := app_context.GetTimezoneFromContext(ctx)
		if user.Timezone != newTimezone {
			user.Timezone = newTimezone
			if err := u.userRepo.UpdateTimezone(txCtx, user.ID, newTimezone); err != nil {
				return err
			}
		}

		user.Tokens = &entities.Tokens{}

		user.Tokens.AccessToken, err = u.tokenService.GenerateJWTToken(map[string]any{"user_id": user.ID}, nil)
		if err != nil {
			return err
		}

		if err := u.refreshTokenRepo.RevokeAllByUserID(txCtx, user.ID); err != nil {
			return err
		}

		refreshToken, err := u.tokenService.GenerateRefreshToken(user.ID, newTimezone)
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
