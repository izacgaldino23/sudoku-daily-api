package domain

import (
	"sudoku-daily-api/src/domain/entities"
)

type (
	SudokuGenerator interface {
		GenerateDaily(size entities.BoardSize, seed int64) (*entities.Sudoku, error)
	}

	PasswordHasher interface {
		Hash(password string) (string, error)
		Compare(password, encodedHash string) error
	}

	TokenService interface {
		GenerateAccessToken(userID string) (string, error)
		GenerateRefreshToken(userID string) (string, error)
		ValidateAccessToken(token string) (string, error)
	}
)
