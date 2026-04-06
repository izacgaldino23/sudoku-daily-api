package pkg

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrNotFound               = newError("not_found")
	ErrQueryParamInvalid      = newError("invalid_query_param")
	ErrInvalidEmail           = newError("invalid_email")
	ErrInvalidToken           = newError("invalid_token")
	ErrTokenExpired           = newError("token_expired")
	ErrInvalidCredentials     = newError("invalid_credentials")
	ErrEmailAlreadyRegistered = newError("email_already_registered")
	ErrRefreshTokenExpired    = newError("refresh_token_expired")
	ErrRefreshTokenRevoked    = newError("refresh_token_revoked")
	ErrBodyInvalid            = newError("invalid_body")
	ErrInvalidSolution        = newError("invalid_solution")
	ErrInvalidLeaderboardType = newError("invalid_leaderboard_type")
	ErrSizeRequired           = newError("size_required")
	ErrSizeNotAllowed         = newError("size_not_allowed")
	ErrInvalidSize            = newError("invalid_size")
	ErrInvalidLimit           = newError("invalid_limit")
	ErrInvalidPage            = newError("invalid_page")
	ErrInternalServerError    = newError("internal_server_error")
	ErrTooManyRequests        = newError("too_many_requests")
	ErrAlreadyPlayed          = newError("already_played")
)

type (
	Error struct {
		Code          string            `json:"code"`
		Message       string            `json:"message"`
		ValidationErr []ValidationError `json:"validation_errors,omitempty"`
	}
)

func (e *Error) Error() string {
	return e.Code + ": " + e.Message
}

func newError(code string) *Error {
	return &Error{Code: code, Message: code}
}

func NewError(message string) *Error {
	return &Error{Code: "internal_server_error", Message: message}
}

func FromError(err error) *Error {
	if validationErrs, ok := err.(ValidationErrors); ok {
		return &Error{
			Code:          "validation_error",
			Message:       validationErrs.Error(),
			ValidationErr: validationErrs,
		}
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{Code: "internal_server_error", Message: err.Error()}
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
	case ErrEmailAlreadyRegistered, ErrQueryParamInvalid, ErrBodyInvalid, ErrSizeRequired, ErrSizeNotAllowed, ErrInvalidSize, ErrInvalidLimit, ErrInvalidPage, ErrInvalidLeaderboardType:
		return http.StatusBadRequest
	case ErrInvalidCredentials, ErrRefreshTokenExpired, ErrRefreshTokenRevoked, ErrInvalidToken, ErrInvalidEmail, ErrTokenExpired:
		return http.StatusUnauthorized
	case ErrTooManyRequests:
		return http.StatusTooManyRequests
	case ErrAlreadyPlayed:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
