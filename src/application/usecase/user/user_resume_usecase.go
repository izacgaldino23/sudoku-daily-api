package user

import (
	"context"

	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
)

type (
	UserResumeUseCase interface {
		Execute(ctx context.Context, userID vo.UUID) (*entities.Resume, error)
	}

	userResumeUseCase struct {
		resumeFetcher domain.ResumeFetcher
	}
)

func NewUserResumeUseCase(resumeFetcher domain.ResumeFetcher) UserResumeUseCase {
	return &userResumeUseCase{resumeFetcher: resumeFetcher}
}

func (u *userResumeUseCase) Execute(ctx context.Context, userID vo.UUID) (*entities.Resume, error) {
	var (
		resume = &entities.Resume{}
		err    error
	)

	// Get my total solved games by size
	resume.TotalGames, err = u.resumeFetcher.GetTotalSolvedByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get my today solved games by size
	resume.TodayGames, err = u.resumeFetcher.GetTodaySolvedByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get my best times by size
	resume.BestTimes, err = u.resumeFetcher.GetBestTimesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return resume, err
}
