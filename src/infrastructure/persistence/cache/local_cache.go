package cache

import (
	"sudoku-daily-api/src/domain"
	"sync"
	"time"
)

type (
	localCache struct {
		sync.RWMutex
		data      map[string]any
		order     []string
		maxSize   int
		updatedAt time.Time
	}
)

func NewLocalCache(maxSize int) domain.Cache {
	return &localCache{
		data:      make(map[string]any),
		order:     make([]string, 0),
		maxSize:   maxSize,
		updatedAt: time.Now(),
	}
}

func (c *localCache) Get(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()

	data, ok := c.data[key]

	return data, ok
}

func (c *localCache) Set(key string, value any) {
	c.Lock()
	defer c.Unlock()

	if len(c.data) < c.maxSize {
		c.order = append(c.order, key)
	} else {
		delete(c.data, c.order[0])
		c.order = c.order[1:]
	}

	c.data[key] = value
	c.order = append(c.order, key)
	c.updatedAt = time.Now()
}

func (c *localCache) Flush() {
	c.Lock()
	defer c.Unlock()

	c.data = make(map[string]any)
	c.order = make([]string, 0)
	c.updatedAt = time.Now()
}
