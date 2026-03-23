package bootstrap

import (
	persistenceRefreshToken "sudoku-daily-api/src/infrastructure/persistence/database/refresh_token"
	persistenceStats "sudoku-daily-api/src/infrastructure/persistence/database/stats"
	persistenceSudoku "sudoku-daily-api/src/infrastructure/persistence/database/sudoku"
	persistenceTx "sudoku-daily-api/src/infrastructure/persistence/database/tx"
	persistenceUser "sudoku-daily-api/src/infrastructure/persistence/database/user"
)

func (c *Container) BuildRepositories() {
	db := c.DB.BunConnection

	c.SudokuRepository = persistenceSudoku.NewRepository(db)
	c.UserRepository = persistenceUser.NewRepository(db)
	c.RefreshTokenRepository = persistenceRefreshToken.NewRepository(db)
	c.UserStatsRepository = persistenceStats.NewRepository(db)

	c.TxManager = persistenceTx.NewTransactionManager(db)
}
