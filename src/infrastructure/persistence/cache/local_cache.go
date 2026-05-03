package cache

import (
	"sync"

	"sudoku-daily-api/src/domain"
)

type (
	localCache struct {
		sync.RWMutex
		data       map[string]any
		order      []string
		orderIndex int
		maxSize    int
	}
)

func NewLocalCache(maxSize int) domain.Cache {
	return &localCache{
		data:    make(map[string]any, maxSize),
		order:   make([]string, maxSize),
		maxSize: maxSize,
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

	if _, exists := c.data[key]; exists {
		c.data[key] = value
		return
	}

	if len(c.data) >= c.maxSize {
		delete(c.data, c.order[c.orderIndex])
	}

	c.order[c.orderIndex] = key
	c.data[key] = value

	c.orderIndex++
	if c.orderIndex >= c.maxSize {
		c.orderIndex = 0
	}
}

func (c *localCache) Flush() {
	c.Lock()
	defer c.Unlock()

	c.data = make(map[string]any, c.maxSize)
	c.order = make([]string, c.maxSize)
	c.orderIndex = 0
}
