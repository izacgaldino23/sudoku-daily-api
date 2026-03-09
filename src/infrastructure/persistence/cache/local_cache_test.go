package cache_test

import (
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/infrastructure/persistence/cache"
	"testing"
)

func TestLocalCache_SetAndGet(t *testing.T) {
	// Definimos a estrutura dos casos de teste
	tests := []struct {
		name          string
		maxSize       int
		actions       func(c domain.Cache)
		keyToFetch    string
		expectedValue any
		expectedFound bool
	}{
		{
			name:    "Should get a value inserted",
			maxSize: 10,
			actions: func(c domain.Cache) {
				c.Set("user_1", "Gabriel")
			},
			keyToFetch:    "user_1",
			expectedValue: "Gabriel",
			expectedFound: true,
		},
		{
			name:    "Should return false if key is not found",
			maxSize: 10,
			actions: func(c domain.Cache) {
				c.Set("user_1", "Gabriel")
			},
			keyToFetch:    "user_2",
			expectedValue: nil,
			expectedFound: false,
		},
		{
			name:    "Should update a previously inserted value",
			maxSize: 10,
			actions: func(c domain.Cache) {
				c.Set("api_key", "12345")
				c.Set("api_key", "67890")
			},
			keyToFetch:    "api_key",
			expectedValue: "67890",
			expectedFound: true,
		},
		{
			name:    "Should respect the max size (LRU/FIFO)",
			maxSize: 2,
			actions: func(c domain.Cache) {
				c.Set("k1", "v1")
				c.Set("k2", "v2")
				c.Set("k3", "v3")
			},
			keyToFetch:    "k1",
			expectedValue: nil,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c := cache.NewLocalCache(tt.maxSize)

			// Act
			tt.actions(c)
			gotValue, gotFound := c.Get(tt.keyToFetch)

			// Assert
			if gotFound != tt.expectedFound {
				t.Errorf("Get() found = %v, want %v", gotFound, tt.expectedFound)
			}
			if gotValue != tt.expectedValue {
				t.Errorf("Get() value = %v, want %v", gotValue, tt.expectedValue)
			}
		})
	}
}
