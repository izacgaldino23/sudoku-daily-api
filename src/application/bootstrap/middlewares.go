package bootstrap

import (
	"sudoku-daily-api/src/infrastructure/http/middlewares"

	"github.com/rs/zerolog/log"
)

func (c *Container) BuildMiddlewares() {
	c.RequireJWT = middlewares.RequireJWTMiddleware(c.TokenService)
	c.OptionalJWT = middlewares.OptionalJWTMiddleware(c.TokenService)
	c.AuthMinimum = middlewares.AuthMinimumMiddleware(c.TokenService)
	c.Session = middlewares.SessionMiddleware(c.TokenService)
	c.LogMiddleware = middlewares.LogMiddleware(log.Logger)
	c.RequestIDMiddleware = middlewares.NewRequestIDMiddleware()
}
