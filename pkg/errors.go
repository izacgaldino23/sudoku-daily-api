package pkg

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrQueryParamInvalid      = newError(ErrorCodeQueryParamInvalid, "invalid query param", http.StatusBadRequest)
	ErrInvalidEmail           = newError(ErrorCodeInvalidEmail, "invalid email", http.StatusUnauthorized)
	ErrInvalidToken           = newError(ErrorCodeInvalidToken, "invalid token", http.StatusUnauthorized)
	ErrTokenExpired           = newError(ErrorCodeTokenExpired, "token expired", http.StatusUnauthorized)
	ErrInvalidCredentials     = newError(ErrorCodeInvalidCredentials, "invalid credentials", http.StatusUnauthorized)
	ErrEmailAlreadyRegistered = newError(ErrorCodeEmailAlreadyRegistered, "email already registered", http.StatusBadRequest)
	ErrRefreshTokenExpired    = newError(ErrorCodeRefreshTokenExpired, "refresh token expired", http.StatusUnauthorized)
	ErrRefreshTokenRevoked    = newError(ErrorCodeRefreshTokenRevoked, "refresh token revoked", http.StatusUnauthorized)
	ErrBodyInvalid            = newError(ErrorCodeBodyInvalid, "invalid body", http.StatusBadRequest)
	ErrInvalidSolution        = newError(ErrorCodeInvalidSolution, "invalid solution", http.StatusBadRequest)
	ErrInvalidLeaderboardType = newError(ErrorCodeInvalidLeaderboardType, "invalid leaderboard type", http.StatusBadRequest)
	ErrInternalServerError    = newError(ErrorCodeInternalServerError, "internal server error", http.StatusInternalServerError)
	ErrTooManyRequests        = newError(ErrorCodeTooManyRequests, "too many requests", http.StatusTooManyRequests)
	ErrAlreadyPlayed          = newError(ErrorCodeAlreadyPlayed, "user has already played", http.StatusConflict)

	ErrUserNotFound         = newError(ErrorCodeUserNotFound, "user not found", http.StatusNotFound)
	ErrSudokuNotFound       = newError(ErrorCodeSudokuNotFound, "sudoku not found", http.StatusNotFound)
	ErrRefreshTokenNotFound = newError(ErrorCodeRefreshTokenNotFound, "refresh token not found", http.StatusNotFound)
	ErrSolutionNotFound     = newError(ErrorCodeSolutionNotFound, "solution not found", http.StatusNotFound)
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
