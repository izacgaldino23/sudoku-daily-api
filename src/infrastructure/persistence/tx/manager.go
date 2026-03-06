package tx

import (
	"context"

	"github.com/uptrace/bun"
)

type (
	txKey struct{}

	Manager struct {
		db *bun.DB
	}
)

func NewManager(db *bun.DB) *Manager {
	return &Manager{db: db}
}

func (m *Manager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := extractTx(ctx); ok {
		return fn(ctx)
	}

	return m.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ctxWithTx := injectTx(ctx, tx)
		return fn(ctxWithTx)
	})
}

func (m *Manager) GetExecutor(ctx context.Context) bun.IDB {
	if tx, ok := extractTx(ctx); ok {
		return tx
	}

	return m.db
}

func injectTx(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) (bun.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(bun.Tx)
	return tx, ok
}
