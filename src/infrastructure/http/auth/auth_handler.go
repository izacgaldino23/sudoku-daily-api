package auth

import (
	"net/http"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/application/usecase/user"

	"github.com/gofiber/fiber/v3"
)

type (
	AuthHandler interface {
		Register(c fiber.Ctx) error
	}

	authHandler struct {
		userRegisterUseCase user.UserRegisterUseCase
	}
)

func NewAuthHandler(
	userRegisterUseCase user.UserRegisterUseCase,
) AuthHandler {
	return &authHandler{
		userRegisterUseCase: userRegisterUseCase,
	}
}

func (a *authHandler) Register(c fiber.Ctx) error {
	var (
		req RegisterRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	user, err := a.userRegisterUseCase.Execute(c.Context(), req.ToDomain())
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.Status(http.StatusOK).JSON(user)
}
