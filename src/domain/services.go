package domain

import (
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
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
		GenerateAccessToken(userID vo.UUID) (string, error)
		GenerateRefreshToken(userID vo.UUID) (*entities.RefreshToken, error)
		ValidateAccessToken(token string) (string, error)
	}
)
