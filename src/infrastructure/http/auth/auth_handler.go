package auth

import (
	"net/http"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/pkg/config"
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
		Resume(c fiber.Ctx) error
	}

	authHandler struct {
		userRegisterUseCase     user.UserRegisterUseCase
		userLoginUseCase        user.UserLoginUseCase
		userRefreshTokenUseCase user.UserRefreshTokenUseCase
		userLogoutUseCase       user.UserLogoutUseCase
		userResumeUseCase       user.UserResumeUseCase
	}
)

func NewAuthHandler(
	userRegisterUseCase user.UserRegisterUseCase,
	userLoginUseCase user.UserLoginUseCase,
	userRefreshTokenUseCase user.UserRefreshTokenUseCase,
	userLogoutUseCase user.UserLogoutUseCase,
	userResumeUseCase user.UserResumeUseCase,
) AuthHandler {
	return &authHandler{
		userRegisterUseCase:     userRegisterUseCase,
		userLoginUseCase:        userLoginUseCase,
		userRefreshTokenUseCase: userRefreshTokenUseCase,
		userLogoutUseCase:       userLogoutUseCase,
		userResumeUseCase:       userResumeUseCase,
	}
}

// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration request"
// @Success 201 {string} string "User created"
// @Failure 400 {object} pkg.Error "invalid_body, validation_error"
// @Failure 409 {object} pkg.Error "email_already_registered"
// @Router /api/auth/register [post]
func (a *authHandler) Register(c fiber.Ctx) error {
	var (
		request RegisterRequest
	)
	if err := c.Bind().Body(&request); err != nil {
		return pkg.ErrBodyInvalid
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return pkg.JsonError(c, err)
	}

	_, err := a.userRegisterUseCase.Execute(c.Context(), request.ToDomain())
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.Status(http.StatusCreated).SendString("")
}

// @Summary Login user
// @Description Authenticates a user and returns access token. Refresh token is sent via secure HTTP-only cookie.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} pkg.Error "invalid_body"
// @Failure 401 {object} pkg.Error "invalid_credentials"
// @Router /api/auth/login [post]
func (a *authHandler) Login(c fiber.Ctx) error {
	var (
		request LoginRequest
	)
	if err := c.Bind().Body(&request); err != nil {
		return pkg.ErrBodyInvalid
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return pkg.JsonError(c, err)
	}

	userData, err := a.userLoginUseCase.Execute(c.Context(), request.ToDomain())
	if err != nil {
		return pkg.JsonError(c, err)
	}

	refreshTokenDuration := config.GetConfig().Auth.RefreshTokenDuration
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    userData.Tokens.RefreshToken,
		Path:    "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   refreshTokenDuration,
	})

	resp := LoginResponse{}
	resp.FromDomain(userData)

	return c.Status(http.StatusOK).JSON(resp)
}

// @Summary Refresh access token
// @Description Refreshes an expired access token using refresh token from cookie. Returns new access token and rotates refresh token cookie.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} RefreshTokenResponse
// @Failure 401 {object} pkg.Error "invalid_token, refresh_token_expired, refresh_token_revoked"
// @Router /api/auth/refresh [post]
func (a *authHandler) Refresh(c fiber.Ctx) error {
	refreshTokenCookie := c.Cookies("refresh_token")
	if refreshTokenCookie == "" {
		return pkg.JsonError(c, pkg.ErrInvalidToken)
	}

	accessToken, newRefreshToken, err := a.userRefreshTokenUseCase.Execute(c.Context(), refreshTokenCookie)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	refreshTokenDuration := config.GetConfig().Auth.RefreshTokenDuration
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:    "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   refreshTokenDuration,
	})

	refreshTokenResponse := RefreshTokenResponse{}
	refreshTokenResponse.AccessToken = accessToken

	return c.Status(http.StatusOK).JSON(refreshTokenResponse)
}

// @Summary Logout user
// @Description Invalidates the user's refresh token from cookie
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "Logged out successfully"
// @Failure 401 {object} pkg.Error "invalid_token"
// @Router /api/auth/logout [post]
func (a *authHandler) Logout(c fiber.Ctx) error {
	var (
		userID  vo.UUID
		reqCtx  = c.Context()
	)

	refreshTokenCookie := c.Cookies("refresh_token")
	if refreshTokenCookie == "" {
		return pkg.JsonError(c, pkg.ErrInvalidToken)
	}

	userID = app_context.GetUserIDFromContext(reqCtx)

	err := a.userLogoutUseCase.Execute(reqCtx, userID, refreshTokenCookie)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:   "refresh_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	return c.SendStatus(http.StatusOK)
}

// @Summary Get user resume
// @Description Returns user statistics including total games, today's games, and best times
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ResumeResponse
// @Router /api/auth/resume [get]
func (a *authHandler) Resume(c fiber.Ctx) error {
	var reqCtx = c.Context()
	userID := app_context.GetUserIDFromContext(reqCtx)

	resume, err := a.userResumeUseCase.Execute(reqCtx, userID)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	response := ResumeResponse{}
	response.FromDomain(resume)

	return c.Status(http.StatusOK).JSON(response)
}
