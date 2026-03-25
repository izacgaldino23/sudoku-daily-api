package domain

import (
	"context"
	"time"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
)

type (
	SudokuGenerator interface {
		GenerateDaily(size entities.BoardSize, date time.Time) (*entities.Sudoku, error)
	}

	PasswordHasher interface {
		Hash(password string) (string, error)
		Compare(password, encodedHash string) error
	}

	TokenService interface {
		GenerateJWTToken(map[string]any, *int) (string, error)
		GenerateRefreshToken(userID vo.UUID) (*entities.RefreshToken, error)
		ValidateAccessToken(token string) (vo.UUID, error)
		ParseToken(token string) (result map[string]any, err error)
	}

	SudokuDailyFetcher interface {
		GetDaily(ctx context.Context, size entities.BoardSize) (*entities.Sudoku, error)
		GetByDateAndSize(ctx context.Context, date time.Time, size entities.BoardSize) (*entities.Sudoku, error)
	}

	ResumeFetcher interface {
		GetTotalSolvedByUser(ctx context.Context, userID vo.UUID) (map[entities.BoardSize]int, error)
		GetTodaySolvedByUser(ctx context.Context, userID vo.UUID) ([]entities.GameResult, error)
		GetBestTimesByUser(ctx context.Context, userID vo.UUID) ([]entities.GameResult, error)
	}
)
