package pkg

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrQueryParamInvalid      = newError("invalid_query_param", "invalid query param", http.StatusBadRequest)
	ErrInvalidEmail           = newError("invalid_email", "invalid email", http.StatusUnauthorized)
	ErrInvalidToken           = newError("invalid_token", "invalid token", http.StatusUnauthorized)
	ErrTokenExpired           = newError("token_expired", "token expired", http.StatusUnauthorized)
	ErrInvalidCredentials     = newError("invalid_credentials", "invalid credentials", http.StatusUnauthorized)
	ErrEmailAlreadyRegistered = newError("email_already_registered", "email already registered", http.StatusBadRequest)
	ErrRefreshTokenExpired    = newError("refresh_token_expired", "refresh token expired", http.StatusUnauthorized)
	ErrRefreshTokenRevoked    = newError("refresh_token_revoked", "refresh token revoked", http.StatusUnauthorized)
	ErrBodyInvalid            = newError("invalid_body", "invalid body", http.StatusBadRequest)
	ErrInvalidSolution        = newError("invalid_solution", "invalid solution", http.StatusBadRequest)
	ErrInvalidLeaderboardType = newError("invalid_leaderboard_type", "invalid leaderboard type", http.StatusBadRequest)
	ErrInternalServerError    = newError("internal_server_error", "internal server error", http.StatusInternalServerError)
	ErrTooManyRequests        = newError("too_many_requests", "too many requests", http.StatusTooManyRequests)
	ErrAlreadyPlayed          = newError("already_played", "user has already played", http.StatusConflict)

	ErrUserNotFound         = newError("user_not_found", "user not found", http.StatusNotFound)
	ErrSudokuNotFound       = newError("sudoku_not_found", "sudoku not found", http.StatusNotFound)
	ErrRefreshTokenNotFound = newError("refresh_token_not_found", "refresh token not found", http.StatusNotFound)
	ErrSolutionNotFound     = newError("solution_not_found", "solution not found", http.StatusNotFound)
)

type (
	Error struct {
		Code          string            `json:"code"`
		Message       string            `json:"message"`
		StatusCode    int               `json:"-"`
		Err           error             `json:"-"`
		ValidationErr []ValidationError `json:"validation_errors,omitempty"`
	}
)

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Message + ": " + e.Err.Error()
	}
	return e.Code + ": " + e.Message
}

func newError(code, message string, statusCode int) *Error {
	return &Error{Code: code, Message: message, StatusCode: statusCode}
}

func NewError(message string) *Error {
	return &Error{Code: "internal_server_error", Message: message, StatusCode: http.StatusInternalServerError}
}

func FromError(err error) *Error {
	if validationErrs, ok := err.(ValidationErrors); ok {
		return &Error{
			Code:          "validation_error",
			Message:       validationErrs.Error(),
			StatusCode:    http.StatusBadRequest,
			ValidationErr: validationErrs,
		}
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{Code: "internal_server_error", Message: err.Error(), StatusCode: http.StatusInternalServerError}
}

func JsonError(c fiber.Ctx, err error) error {
	appErr := FromError(err)
	if appErr.StatusCode == 0 {
		appErr.StatusCode = http.StatusInternalServerError
	}

	return c.Status(appErr.StatusCode).JSON(appErr)
}

func JsonErrorWithStatus(c fiber.Ctx, err error, status int) error {
	appErr := FromError(err)
	appErr.StatusCode = status
	return c.Status(status).JSON(appErr)
}
