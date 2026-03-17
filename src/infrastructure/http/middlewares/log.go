package middlewares

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/app_context"
)

func LogMiddleware(base zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		reqCtx := c.Context()
		requestID := app_context.GetRequestIDFromContext(reqCtx)

		log := base.With().
			Str("path", c.Path()).
			Str("method", c.Method()).
			Str("request_id", requestID.String()).
			Logger()

		ctx := log.WithContext(reqCtx)

		c.SetContext(ctx)

		err := c.Next()
		if err != nil {
			log.Error().Err(err)
		} else if c.Response().StatusCode() >= 300 {
			var reqErr pkg.Error
			errResponseData := c.Response().Body()

			err = json.Unmarshal(errResponseData, &reqErr)
			if err != nil {
				log.Error().Msg(err.Error())
			} else {
				log.Error().Msg(reqErr.Message)
			}
		}

		return err
	}
}
