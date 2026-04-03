package bootstrap

import (
	"sudoku-daily-api/src/infrastructure/http/middlewares"

	"github.com/rs/zerolog/log"
)

func (c *Container) BuildMiddlewares() {
	c.Middlewares.RequireJWT = middlewares.RequireJWTMiddleware(c.TokenService)
	c.Middlewares.OptionalJWT = middlewares.OptionalJWTMiddleware(c.TokenService)
	c.Middlewares.AuthMinimum = middlewares.AuthMinimumMiddleware(c.TokenService)
	c.Middlewares.Session = middlewares.SessionMiddleware(c.TokenService)
	c.Middlewares.LogMiddleware = middlewares.LogMiddleware(log.Logger)
	c.Middlewares.RequestID = middlewares.NewRequestIDMiddleware()
	c.Middlewares.ResponseHeaders = middlewares.NewResponseHeadersMiddleware()
	c.Middlewares.AuthOIDC = middlewares.AuthOIDCMiddleware(c.Config.Auth.OidcEnabled, c.Config.Auth.OidcAudience)
	c.Middlewares.GlobalRateLimiter = middlewares.NewGlobalRateLimiterMiddleware(c.Config.Limits.MaxRequestCountGlobal)
	c.Middlewares.UserRateLimiter = middlewares.NewUserRateLimitMiddleware(c.Config.Limits.MaxRequestCountUser)
}
