package leaderboard_usecase

import (
	"context"
	"strconv"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
)

type (
	GetLeaderboard interface {
		Execute(ctx context.Context, params *entities.LeaderboardSearchParams) (*entities.Leaderboard, error)
	}

	leaderboardUsecase struct {
		userStatsRepository repository.UserStatsRepository
		sudokuRepository    repository.SudokuRepository
		sudokuFetcher       domain.SudokuDailyFetcher
	}
)

func NewLeaderboardUsecase(
	userStatsRepository repository.UserStatsRepository,
	sudokuRepository repository.SudokuRepository,
	sudokuFetcher domain.SudokuDailyFetcher,
) GetLeaderboard {
	return &leaderboardUsecase{
		userStatsRepository: userStatsRepository,
		sudokuRepository:    sudokuRepository,
		sudokuFetcher:       sudokuFetcher,
	}
}

func (l *leaderboardUsecase) Execute(ctx context.Context, params *entities.LeaderboardSearchParams) (*entities.Leaderboard, error) {
	switch params.Type {
	case entities.DailyLeaderboardType.String():
		return l.getDaily(ctx, params)
	case entities.AllTimeLeaderboardType.String():
		return l.getAllTimeBest(ctx, params)
	case entities.TotalSolvesLeaderboardType.String():
		return l.getByTotalSolves(ctx, params)
	case entities.StreakLeaderboardType.String():
		return l.getByStreak(ctx, params)
	default:
		return nil, pkg.ErrInvalidLeaderboardType
	}
}

func (l *leaderboardUsecase) getDaily(ctx context.Context, params *entities.LeaderboardSearchParams) (*entities.Leaderboard, error) {
	sudoku, err := l.sudokuFetcher.GetDaily(ctx, params.Size)
	if err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.Limit

	solves, hasNext, err := l.sudokuRepository.GetDailyLeaderboard(ctx, sudoku.ID, params.Limit, offset)
	if err != nil {
		return nil, err
	}

	return l.solvesToLeaderboardEntries(solves, hasNext), nil
}

func (l *leaderboardUsecase) getAllTimeBest(ctx context.Context, params *entities.LeaderboardSearchParams) (*entities.Leaderboard, error) {
	offset := (params.Page - 1) * params.Limit

	solves, hasNext, err := l.sudokuRepository.GetAllTimeBestLeaderboard(ctx, params.Size, params.Limit, offset)
	if err != nil {
		return nil, err
	}

	return l.solvesToLeaderboardEntries(solves, hasNext), nil
}

func (l *leaderboardUsecase) solvesToLeaderboardEntries(solves []entities.Solve, hasNext bool) *entities.Leaderboard {
	leaderboard := &entities.Leaderboard{
		HasNext: hasNext,
		Entries: make([]entities.Entry, len(solves)),
	}

	for i, solve := range solves {
		leaderboard.Entries[i] = entities.Entry{
			Rank:     i + 1,
			Username: solve.Username,
			Value:    strconv.FormatInt(int64(solve.Duration), 10),
		}
	}

	return leaderboard
}

func (l *leaderboardUsecase) getByTotalSolves(ctx context.Context, params *entities.LeaderboardSearchParams) (*entities.Leaderboard, error) {
	offset := (params.Page - 1) * params.Limit

	stats, hasNext, err := l.userStatsRepository.GetTotalSolvesLeaderboard(ctx, params.Limit, offset)
	if err != nil {
		return nil, err
	}

	return l.statsToLeaderboardEntries(stats, hasNext), nil
}

func (l *leaderboardUsecase) getByStreak(ctx context.Context, params *entities.LeaderboardSearchParams) (*entities.Leaderboard, error) {
	offset := (params.Page - 1) * params.Limit

	stats, hasNext, err := l.userStatsRepository.GetBestStreakLeaderboard(ctx, params.Limit, offset, time.Now())
	if err != nil {
		return nil, err
	}

	return l.statsToLeaderboardEntries(stats, hasNext), nil
}

func (l *leaderboardUsecase) statsToLeaderboardEntries(solves []entities.UserStats, hasNext bool) *entities.Leaderboard {
	leaderboard := &entities.Leaderboard{
		HasNext: hasNext,
		Entries: make([]entities.Entry, len(solves)),
	}

	for i, stat := range solves {
		leaderboard.Entries[i] = entities.Entry{
			Rank:     i + 1,
			Username: stat.Username,
			Value:    strconv.FormatInt(int64(stat.TotalSolved), 10),
		}
	}

	return leaderboard
}
