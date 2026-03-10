package bootstrap

import (
	sudokuUsecase "sudoku-daily-api/src/application/usecase/sudoku"
	userUsecase "sudoku-daily-api/src/application/usecase/user"
)

func (c *Container) BuildUseCases() {

	c.GetDailySudoku = sudokuUsecase.NewSudokuGetDailyUseCase(
		c.TokenService,
		c.SudokuFetcher,
	)

	c.GenerateAllSudokus = sudokuUsecase.NewSudokuGenerateAllUseCase(
		c.TxManager,
		c.SudokuRepository,
		c.GeneratorService,
	)

	c.VerifySolution = sudokuUsecase.NewSudokuVerifySolutionUseCase(
		c.UserRepository,
		c.SudokuRepository,
		c.TokenService,
		c.SudokuFetcher,
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
}