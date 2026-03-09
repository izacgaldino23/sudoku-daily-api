package app_context

import (
	"context"
	"sudoku-daily-api/src/domain/vo"
)

type contextUserIDKey struct{}

func SetUserOnContext(ctx context.Context, userID vo.UUID) context.Context {
	newCtx := context.WithValue(ctx, contextUserIDKey{}, userID)
	return newCtx
}

func GetUserIDFromContext(ctx context.Context) vo.UUID {
	return ctx.Value(contextUserIDKey{}).(vo.UUID)
}
