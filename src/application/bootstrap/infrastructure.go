package bootstrap

import (
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/infrastructure/persistence/cache"
)

const maxCacheSize = 50

func (c *Container) BuildInfrastructure() {
	c.Config = config.GetConfig()
	c.DB = database.GetDB()
	c.LocalCache = cache.NewLocalCache(maxCacheSize)
}
