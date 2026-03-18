package bootstrap

import (
	"sudoku-daily-api/src/domain/strategies"
	"sudoku-daily-api/src/services"
)

func (c *Container) BuildServices() {

	authConfig := c.Config.Auth

	fillStrategy := strategies.NewFillStrategy()
	hideStrategy := strategies.NewHideStrategy()

	c.GeneratorService = services.NewGenerator(fillStrategy, hideStrategy)

	c.PasswordHasher = services.NewPasswordHasher(
		authConfig.Iterations,
		authConfig.Memory,
		authConfig.Parallelism,
		authConfig.KeyLen,
		authConfig.SaltLen,
	)

	c.TokenService = services.NewTokenService(
		authConfig.SecretKey,
		authConfig.AccessTokenDuration,
		authConfig.RefreshTokenDuration,
	)

	c.SudokuFetcher = services.NewSudokuDailyFetcher(
		c.LocalCache,
		c.SudokuRepository,
	)

	c.ResumeFetcher = services.NewResumeFetcher(
		c.SudokuRepository,
	)
}
