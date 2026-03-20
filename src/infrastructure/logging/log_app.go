package logging

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Log(ctx context.Context) *zerolog.Logger {
	return log.Ctx(ctx)

}
