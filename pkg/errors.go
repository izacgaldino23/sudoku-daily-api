package pkg

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrNotFound               = NewError("not found")
	ErrQueryParamInvalid      = NewError("invalid query param")
	ErrInvalidEmail           = NewError("invalid email")
	ErrInvalidToken           = NewError("invalid token")
	ErrTokenExpired           = NewError("token expired")
	ErrInvalidCredentials     = NewError("invalid credentials")
	ErrEmailAlreadyRegistered = NewError("email already registered")
	ErrRefreshTokenExpired    = NewError("refresh token expired")
	ErrRefreshTokenRevoked    = NewError("refresh token revoked")
	ErrBodyInvalid            = NewError("invalid body")
	ErrInvalidSolution        = NewError("invalid solution")
	ErrInvalidLeaderboardType = NewError("invalid leaderboard type")
	ErrInternalServerError    = NewError("internal server error")
)

type (
	Error struct {
		Message       string            `json:"message"`
		ValidationErr []ValidationError `json:"validation_errors,omitempty"`
	}
)

func (e *Error) Error() string {
	return e.Message
}

func FromError(err error) *Error {
	if validationErrs, ok := err.(ValidationErrors); ok {
		return &Error{
			Message:       validationErrs.Error(),
			ValidationErr: validationErrs,
		}
	}
	return &Error{Message: err.Error()}
}

func NewError(message string) *Error {
	return &Error{Message: message}
}

func JsonError(c fiber.Ctx, err error) error {
	status := MapErrorToStatus(err)
	err = FromError(err)

	return c.Status(status).JSON(err)
}

func JsonErrorWithStatus(c fiber.Ctx, err error, status int) error {
	return c.Status(status).JSON(FromError(err))
}

func MapErrorToStatus(err error) int {
	if _, ok := err.(ValidationErrors); ok {
		return http.StatusBadRequest
	}

	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrEmailAlreadyRegistered, ErrQueryParamInvalid, ErrBodyInvalid:
		return http.StatusBadRequest
	case ErrInvalidCredentials, ErrRefreshTokenExpired, ErrRefreshTokenRevoked, ErrInvalidToken, ErrInvalidEmail, ErrTokenExpired:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
