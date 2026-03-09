package logging

import (
	"context"

	"github.com/rs/zerolog"
)

type ctxKey struct{}

var LoggerKey = ctxKey{}

func Log(ctx context.Context) *zerolog.Logger {
	if log, ok := ctx.Value(LoggerKey).(*zerolog.Logger); ok {
		return log
	}

	l := zerolog.Nop()
	return &l
}
