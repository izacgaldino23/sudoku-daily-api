package domain

import (
	"context"
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
		GenerateJWTToken(map[string]any) (string, error)
		GenerateRefreshToken(userID vo.UUID) (*entities.RefreshToken, error)
		ValidateAccessToken(token string) (vo.UUID, error)
		ParseToken(token string) (result map[string]any, err error)
	}

	SudokuDailyFetcher interface {
		GetDaily(ctx context.Context, size int) (*entities.Sudoku, error)
	}
)
