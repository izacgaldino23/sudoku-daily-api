package bootstrap

import (
	leaderboard_usecase "sudoku-daily-api/src/application/usecase/leaderboard"
	sudokuUsecase "sudoku-daily-api/src/application/usecase/sudoku"
	userUsecase "sudoku-daily-api/src/application/usecase/user"
	user_stats_usecase "sudoku-daily-api/src/application/usecase/user_stats"
)

func (c *Container) BuildUseCases() {

	c.GetDailySudoku = sudokuUsecase.NewSudokuGetDailyUseCase(
		c.TokenService,
		c.SudokuFetcher,
	)

	c.GenerateDailySudokus = sudokuUsecase.NewSudokuGenerateDailyUseCase(
		c.TxManager,
		c.SudokuRepository,
		c.GeneratorService,
		c.SudokuFetcher,
	)

	c.VerifySolution = sudokuUsecase.NewSudokuVerifySolutionUseCase(
		c.SudokuRepository,
		c.TokenService,
		c.SudokuFetcher,
		c.UserStatsSolveAddStrike,
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
		c.RefreshTokenRepository,
		c.TokenService,
	)

	c.UserLogout = userUsecase.NewUserLogoutUseCase(
		c.RefreshTokenRepository,
	)

	c.UserResume = userUsecase.NewUserResumeUseCase(
		c.ResumeFetcher,
	)

	c.UserStatsSolveAddStrike = user_stats_usecase.NewSolveAddStrikeUseCase(
		c.UserStatsRepository,
	)

	c.GetLeaderboardUseCase = leaderboard_usecase.NewLeaderboardUsecase(
		c.UserStatsRepository,
		c.SudokuRepository,
		c.SudokuFetcher,
	)

	c.ResetStrikesUseCase = leaderboard_usecase.NewResetStrikesUseCase(
		c.UserStatsRepository,
	)
}
