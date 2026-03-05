package persistence

import (
	"context"
	"sudoku-daily-api/src/domain/repository"

	"github.com/uptrace/bun"
)

type (
	txKey struct{}

	transactionManager struct{
		db *bun.DB
	}
)

func NewTransactionManager(db *bun.DB) repository.TransactionManager {
	return &transactionManager{
		db: db,
	}
}

func (tm *transactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := extractTx(ctx); ok {
		return fn(ctx)
	}

    return tm.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
        ctxWithTx := injectTx(ctx, tx)
        return fn(ctxWithTx)
    })
}

func (tm *transactionManager) GetExecutor(ctx context.Context) bun.IDB {
	if tx, ok := extractTx(ctx); ok {
		return tx
	}

	return tm.db
}

func injectTx(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) (bun.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(bun.Tx)
	return tx, ok
}
