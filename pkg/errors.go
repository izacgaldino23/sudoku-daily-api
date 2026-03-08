package pkg

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrNotFound               = errors.New("not found")
	ErrQueryParamInvalid      = errors.New("invalid query param")
	ErrInvalidEmail           = errors.New("invalid email")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrRefreshTokenExpired    = errors.New("refresh token expired")
	ErrRefreshTokenRevoked    = errors.New("refresh token revoked")
	ErrBodyInvalid            = errors.New("invalid body")
)

type (
	Error struct {
		Message string `json:"message"`
	}
)

func (e *Error) Error() string {
	return e.Message
}

func FromError(err error) *Error {
	return &Error{Message: err.Error()}
}

func NewError(message string) *Error {
	return &Error{Message: message}
}

func JsonError(c fiber.Ctx, err error) error {
	return c.Status(MapErrorToStatus(err)).JSON(FromError(err))
}

func JsonErrorWithStatus(c fiber.Ctx, err error, status int) error {
	return c.Status(status).JSON(FromError(err))
}

func MapErrorToStatus(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidEmail, ErrEmailAlreadyRegistered, ErrQueryParamInvalid, ErrBodyInvalid:
		return http.StatusBadRequest
	case ErrInvalidCredentials, ErrRefreshTokenExpired, ErrRefreshTokenRevoked:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
