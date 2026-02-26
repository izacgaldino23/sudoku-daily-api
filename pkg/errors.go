package pkg

import "github.com/gofiber/fiber/v3"

type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(message string) *Error {
	return &Error{Message: message}
}

func JsonError(c fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(NewError(msg))
}

func JsonErrorWithStatus(c fiber.Ctx, msg string, status int) error {
	return c.Status(status).JSON(NewError(msg))
}