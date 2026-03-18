package bootstrap

import (
	httpSudoku "sudoku-daily-api/src/infrastructure/http/sudoku"
	"sudoku-daily-api/src/infrastructure/http/auth"
)

func (c *Container) BuildHandlers() {

	c.SudokuHandler = httpSudoku.NewSudokuHandler(
		c.GetDailySudoku,
		c.GenerateAllSudokus,
		c.VerifySolution,
	)

	c.AuthHandler = auth.NewAuthHandler(
		c.UserRegister,
		c.UserLogin,
		c.UserRefreshToken,
		c.UserLogout,
		c.UserResume,
	)
}