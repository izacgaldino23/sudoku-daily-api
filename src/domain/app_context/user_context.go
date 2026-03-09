package app_context

import (
	"context"
	"sudoku-daily-api/src/domain/vo"
)

type contextUserIDKey struct{}

func UserContextKey() interface{} {
	return contextUserIDKey{}
}

func SetUserOnContext(ctx context.Context, userID vo.UUID) context.Context {
	return context.WithValue(ctx, UserContextKey(), userID)
}

func GetUserIDFromContext(ctx context.Context) vo.UUID {
	userID, ok := ctx.Value(UserContextKey()).(vo.UUID)
	if !ok {
		return ""
	}

	return userID
}
