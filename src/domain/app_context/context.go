package app_context

import (
	"context"

	"sudoku-daily-api/src/domain/vo"
)

type (
	contextUserIDKey struct{}
	sessionIDKey     struct{}
	requestIDKey     struct{}
	timezoneKey      struct{}
)

func SetUserIDOnContext(ctx context.Context, userID vo.UUID) context.Context {
	return context.WithValue(ctx, contextUserIDKey{}, userID)
}

func SetSessionIDOnContext(ctx context.Context, sessionID vo.UUID) context.Context {
	return context.WithValue(ctx, sessionIDKey{}, sessionID)
}

func SetRequestIDOnContext(ctx context.Context, requestID vo.UUID) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

func SetTimezoneOnContext(ctx context.Context, timezone string) context.Context {
	return context.WithValue(ctx, timezoneKey{}, timezone)
}

func GetUserIDFromContext(ctx context.Context) vo.UUID {
	userID, ok := ctx.Value(contextUserIDKey{}).(vo.UUID)
	if !ok {
		return ""
	}

	return userID
}

func GetSessionIDFromContext(ctx context.Context) vo.UUID {
	sessionID, ok := ctx.Value(sessionIDKey{}).(vo.UUID)
	if !ok {
		return ""
	}

	return sessionID
}

func GetRequestIDFromContext(ctx context.Context) vo.UUID {
	requestID, ok := ctx.Value(requestIDKey{}).(vo.UUID)
	if !ok {
		return ""
	}

	return requestID
}

func GetTimezoneFromContext(ctx context.Context) string {
	timezone, ok := ctx.Value(timezoneKey{}).(string)
	if !ok {
		return "UTC"
	}

	return timezone
}
