package app_context

import (
	"context"
	"sudoku-daily-api/src/domain/vo"
)

type (
	contextUserIDKey struct{}
	sessionIDKey struct{}
)

func SetUserOnContext(ctx context.Context, userID vo.UUID) context.Context {
	return context.WithValue(ctx, contextUserIDKey{}, userID)
}

func GetUserIDFromContext(ctx context.Context) vo.UUID {
	userID, ok := ctx.Value(contextUserIDKey{}).(vo.UUID)
	if !ok {
		return ""
	}

	return userID
}

func SetSessionIDOnContext(ctx context.Context, sessionID vo.UUID) context.Context {
	return context.WithValue(ctx, sessionIDKey{}, sessionID)
}

func GetSessionIDFromContext(ctx context.Context) vo.UUID {
	sessionID, ok := ctx.Value(sessionIDKey{}).(vo.UUID)
	if !ok {
		return ""
	}

	return sessionID
}