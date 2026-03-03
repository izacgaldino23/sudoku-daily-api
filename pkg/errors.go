package pkg

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrorNotFound     = errors.New("not found")
	QueryParamInvalid = errors.New("invalid query param")
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

func JsonErrorWithStatus(c fiber.Ctx, msg string, status int) error {
	return c.Status(status).JSON(NewError(msg))
}

func MapErrorToStatus(err error) int {
	switch err {
	case ErrorNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
