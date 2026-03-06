package tx

import (
	"context"
	"sudoku-daily-api/src/domain/repository"

	"github.com/uptrace/bun"
)

func NewTransactionManager(db *bun.DB) repository.TransactionManager {
	return &transactionManagerAdapter{
		manager: NewManager(db),
	}
}

type transactionManagerAdapter struct {
	manager *Manager
}

func (a *transactionManagerAdapter) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return a.manager.WithinTransaction(ctx, fn)
}
