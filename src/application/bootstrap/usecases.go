package bootstrap

import (
	leaderboardUsecase "sudoku-daily-api/src/application/usecase/leaderboard"
	sudokuUsecase "sudoku-daily-api/src/application/usecase/sudoku"
	userUsecase "sudoku-daily-api/src/application/usecase/user"
	userStatsUsecase "sudoku-daily-api/src/application/usecase/user_stats"
)

func (c *Container) BuildUseCases() {

	c.GetDailySudoku = sudokuUsecase.NewSudokuGetDailyUseCase(
		c.TokenService,
		c.SudokuFetcher,
		c.SudokuRepository,
	)

	c.GetDailySudokuForGuest = sudokuUsecase.NewSudokuGetDailyForGuestUseCase(
		c.TokenService,
		c.SudokuFetcher,
	)

	c.GenerateDailySudoku = sudokuUsecase.NewSudokuGenerateDailyUseCase(
		c.TxManager,
		c.SudokuRepository,
		c.GeneratorService,
		c.SudokuFetcher,
	)

	c.RemoveUnfinishedAttempts = sudokuUsecase.NewRemoveUnfinishedAttemptsUseCase(
		c.SudokuRepository,
	)

	c.UserStatsSolveAddStrike = userStatsUsecase.NewSolveAddStrikeUseCase(
		c.UserStatsRepository,
	)

	c.VerifySolution = sudokuUsecase.NewSudokuVerifySolutionUseCase(
		c.SudokuRepository,
		c.TokenService,
		c.SudokuFetcher,
		c.UserStatsSolveAddStrike,
		c.TxManager,
	)

	c.VerifySolutionGuest = sudokuUsecase.NewSudokuVerifySolutionGuestUseCase(
		c.SudokuRepository,
		c.TokenService,
		c.SudokuFetcher,
		c.TxManager,
	)

	c.GetUserSolvesUseCase = sudokuUsecase.NewSudokuGetUserSolvesUseCase(
		c.SudokuRepository,
	)

	c.UserRegister = userUsecase.NewUserRegisterUseCase(
		c.UserRepository,
		c.PasswordHasher,
	)

	c.UserLogin = userUsecase.NewUserLoginUseCase(
		c.TxManager,
		c.UserRepository,
		c.RefreshTokenRepository,
		c.PasswordHasher,
		c.TokenService,
	)

	c.UserRefreshToken = userUsecase.NewUserRefreshTokenUseCase(
		c.TxManager,
		c.RefreshTokenRepository,
		c.TokenService,
		c.UserRepository,
	)

	c.UserLogout = userUsecase.NewUserLogoutUseCase(
		c.RefreshTokenRepository,
	)

	c.UserResume = userUsecase.NewUserResumeUseCase(
		c.ResumeFetcher,
	)

	c.GetLeaderboardUseCase = leaderboardUsecase.NewLeaderboardUseCase(
		c.UserStatsRepository,
		c.SudokuRepository,
		c.SudokuFetcher,
	)

	c.ResetStrikesUseCase = leaderboardUsecase.NewResetStrikesUseCase(
		c.UserStatsRepository,
	)
}
