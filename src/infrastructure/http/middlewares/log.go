package middlewares

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/app_context"
)

func severityFromLevel(level zerolog.Level) string {
	switch level {
	case zerolog.DebugLevel:
		return "DEBUG"
	case zerolog.InfoLevel:
		return "INFO"
	case zerolog.WarnLevel:
		return "WARNING"
	case zerolog.ErrorLevel:
		return "ERROR"
	case zerolog.FatalLevel:
		return "CRITICAL"
	case zerolog.PanicLevel:
		return "EMERGENCY"
	default:
		return "DEFAULT"
	}
}

func httpRequestDict(c fiber.Ctx, status int, latency time.Duration) *zerolog.Event {
	return zerolog.Dict().
		Str("requestMethod", c.Method()).
		Str("requestUrl", c.OriginalURL()).
		Int("status", status).
		Str("latency", latency.Round(time.Microsecond).String()).
		Str("userAgent", c.Get("User-Agent")).
		Str("remoteIp", c.IP()).
		Str("protocol", c.Protocol())
}

const serviceName = "sudoku-api"

func serviceContextDict() *zerolog.Event {
	return zerolog.Dict().
		Str("service", serviceName)
}

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
					Str("severity", "INFO").
					RawJSON("body", c.Body()).
					Msg("request received")
			} else {
				log.Info().
					Str("severity", "INFO").
					Bytes("body", c.Body()).
					Msg("request received")
			}
		} else {
			log.Info().
				Str("severity", "INFO").
				Msg("request received")
		}

		c.SetContext(log.WithContext(reqCtx))

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		defer RecordMetrics(c.Method(), c.Path(), status, latency)

		baseEvent := log.With().
			Int("status", status).
			Str("latency", latency.Round(time.Microsecond).String()).
			Dict("httpRequest", httpRequestDict(c, status, latency)).
			Logger()

		if err != nil {
			baseEvent.Error().
				Str("severity", "ERROR").
				Err(err).
				Dict("serviceContext", serviceContextDict()).
				Str("@type", "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent").
				Msg("request failed")
			return err
		}

		if status >= 400 {
			var reqErr pkg.Error
			body := c.Response().Body()

			if len(body) > 0 && json.Unmarshal(body, &reqErr) == nil {
				baseEvent.Error().
					Str("severity", "ERROR").
					Dict("serviceContext", serviceContextDict()).
					Str("@type", "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent").
					Str("error_message", reqErr.Message).
					Msg("http error")
			} else {
				baseEvent.Error().
					Str("severity", "ERROR").
					Bytes("response_body", body).
					Msg("http error (unparsed)")
			}

			return nil
		}

		baseEvent.Info().
			Str("severity", "INFO").
			Msg("request completed")

		return nil
	}
}
