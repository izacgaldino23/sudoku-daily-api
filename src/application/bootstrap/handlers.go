package bootstrap

import (
	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/src/infrastructure/http/leaderboard"
	httpSudoku "sudoku-daily-api/src/infrastructure/http/sudoku"
)

func (c *Container) BuildHandlers() {

	c.SudokuHandler = httpSudoku.NewSudokuHandler(
		c.GetDailySudoku,
		c.GenerateDailySudoku,
		c.VerifySolution,
		c.VerifySolutionGuest,
		c.GetUserSolvesUseCase,
		c.GetDailySudokuForGuest,
		c.RemoveUnfinishedAttempts,
	)

	c.AuthHandler = auth.NewAuthHandler(
		c.UserRegister,
		c.UserLogin,
		c.UserRefreshToken,
		c.UserLogout,
		c.UserResume,
	)

	c.LeaderboardHandler = leaderboard.NewLeaderboardHandler(
		c.GetLeaderboardUseCase,
		c.ResetStrikesUseCase,
	)
}
