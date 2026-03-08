package auth

import (
	"net/http"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/application/usecase/user"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/vo"

	"github.com/gofiber/fiber/v3"
)

type (
	AuthHandler interface {
		Register(c fiber.Ctx) error
		Login(c fiber.Ctx) error
		Refresh(c fiber.Ctx) error
		Logout(c fiber.Ctx) error
	}

	authHandler struct {
		userRegisterUseCase     user.UserRegisterUseCase
		userLoginUseCase        user.UserLoginUseCase
		userRefreshTokenUseCase user.UserRefreshTokenUseCase
		userLogoutUseCase       user.UserLogoutUseCase
	}
)

func NewAuthHandler(
	userRegisterUseCase user.UserRegisterUseCase,
	userLoginUseCase user.UserLoginUseCase,
	userRefreshTokenUseCase user.UserRefreshTokenUseCase,
	userLogoutUseCase user.UserLogoutUseCase,
) AuthHandler {
	return &authHandler{
		userRegisterUseCase:     userRegisterUseCase,
		userLoginUseCase:        userLoginUseCase,
		userRefreshTokenUseCase: userRefreshTokenUseCase,
		userLogoutUseCase:       userLogoutUseCase,
	}
}

func (a *authHandler) Register(c fiber.Ctx) error {
	var (
		req RegisterRequest
	)
	if err := c.Bind().Body(&req); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	if err := pkg.ValidateStruct(req); err != nil {
		return pkg.JsonError(c, err)
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

	if err := pkg.ValidateStruct(req); err != nil {
		return pkg.JsonError(c, err)
	}

	userData, err := a.userLoginUseCase.Execute(c.Context(), req.ToDomain())
	if err != nil {
		return pkg.JsonError(c, err)
	}

	resp := LoginResponse{}
	resp.FromDomain(userData)

	return c.Status(http.StatusOK).JSON(resp)
}

func (a *authHandler) Refresh(c fiber.Ctx) error {
	var (
		req    RefreshTokenRequest
		userID vo.UUID
	)

	if err := c.Bind().Body(&req); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	if err := pkg.ValidateStruct(req); err != nil {
		return pkg.JsonError(c, err)
	}

	userID = app_context.GetUserIDFromContext(c.Context())

	accessToken, err := a.userRefreshTokenUseCase.Execute(c.Context(), req.RefreshToken, userID)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	refreshTokenResponse := RefreshTokenResponse{}
	refreshTokenResponse.AccessToken = accessToken

	return c.Status(http.StatusOK).JSON(refreshTokenResponse)
}

func (a *authHandler) Logout(c fiber.Ctx) error {
	var (
		userID  vo.UUID
		request LogoutRequest
	)

	if err := c.Bind().Body(&request); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	userID = app_context.GetUserIDFromContext(c.Context())

	err := a.userLogoutUseCase.Execute(c.Context(), userID, request.RefreshToken)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.SendStatus(http.StatusOK)
}
