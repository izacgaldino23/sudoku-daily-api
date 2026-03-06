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
		Login(c fiber.Ctx) error
	}

	authHandler struct {
		userRegisterUseCase user.UserRegisterUseCase
		userLoginUseCase    user.UserLoginUseCase
	}
)

func NewAuthHandler(
	userRegisterUseCase user.UserRegisterUseCase,
	userLoginUseCase user.UserLoginUseCase,
) AuthHandler {
	return &authHandler{
		userRegisterUseCase: userRegisterUseCase,
		userLoginUseCase:    userLoginUseCase,
	}
}

func (a *authHandler) Register(c fiber.Ctx) error {
	var (
		req RegisterRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	_, err := a.userRegisterUseCase.Execute(c.Context(), req.ToDomain())
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.SendStatus(http.StatusCreated)
}

func (a *authHandler) Login(c fiber.Ctx) error {
	var (
		req LoginRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	userData, err := a.userLoginUseCase.Execute(c.Context(), req.ToDomain())
	if err != nil {
		return pkg.JsonError(c, err)
	}

	resp := LoginResponse{}
	resp.FromDomain(userData)

	return c.Status(http.StatusOK).JSON(resp)
}
