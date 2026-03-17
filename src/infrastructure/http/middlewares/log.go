package middlewares

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/app_context"
)

func LogMiddleware(base zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		reqCtx := c.Context()
		requestID := app_context.GetRequestIDFromContext(reqCtx)

		log := base.With().
			Str("request_id", requestID.String()).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Logger()

		if c.Method() != fiber.MethodGet && len(c.Body()) < 2048 {
			if json.Valid(c.Body()) {
				log.Info().
					RawJSON("body", c.Body()).
					Msg("request received")
			} else {
				log.Info().
					Bytes("body", c.Body()).
					Msg("request received")
			}
		} else {
			log.Info().Msg("request received")
		}

		c.SetContext(log.WithContext(reqCtx))

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		event := log.With().
			Int("status", status).
			Dur("latency", latency).
			Logger()

		if err != nil {
			event.Error().
				Err(err).
				Msg("request failed")
			return err
		}

		if status >= 400 {
			var reqErr pkg.Error
			body := c.Response().Body()

			if len(body) > 0 && json.Unmarshal(body, &reqErr) == nil {
				event.Error().
					Str("error_message", reqErr.Message).
					Msg("http error")
			} else {
				event.Error().
					Bytes("response_body", body).
					Msg("http error (unparsed)")
			}

			return nil
		}

		event.Info().Msg("request completed")

		return nil
	}
}
