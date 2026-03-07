package app_context

import (
	"context"
	"sudoku-daily-api/src/domain/vo"
)

var contextUserIDKey = struct{}{}

func SetUserOnContext(ctx context.Context, userID vo.UUID) context.Context {
	return context.WithValue(ctx, contextUserIDKey, userID)
}

func GetUserIDFromContext(ctx context.Context) vo.UUID {
	return ctx.Value(contextUserIDKey).(vo.UUID)
}
