package user

import (
	"context"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
)

type (
	UserRegisterUseCase interface {
		Execute(context.Context, *entities.User) (*entities.User, error)
	}

	userRegisterUseCase struct {
		userRepo       repository.UserRepository
		passwordHasher domain.PasswordHasher
	}
)

func NewUserRegisterUseCase(
	userRepo repository.UserRepository,
	passwordHasher domain.PasswordHasher,
) UserRegisterUseCase {
	return &userRegisterUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

// Execute implements [UserRegisterUseCase].
func (u *userRegisterUseCase) Execute(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Validate email
	if !user.Email.IsValid() {
		return nil, pkg.ErrInvalidEmail
	}

	// verify email is not already registered
	existingUser, err := u.userRepo.GetByEmail(ctx, user.Email.String())
	if err != nil && err != pkg.ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, pkg.ErrEmailAlreadyRegistered
	}

	// Hash
	passHash, err := u.passwordHasher.Hash(user.PasswordHash)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = passHash

	user.ID = vo.NewUUID()

	// Create user in database
	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
