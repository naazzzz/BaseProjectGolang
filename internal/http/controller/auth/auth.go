package auth

import (
	"BaseProjectGolang/internal/domain/token"
	"BaseProjectGolang/internal/domain/user"
	"BaseProjectGolang/internal/usecase/user/auth/login"
	"BaseProjectGolang/internal/validation"
	"time"

	"BaseProjectGolang/internal/http/controller"
	"BaseProjectGolang/internal/http/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rotisserie/eris"
)

type Controller struct {
	*controller.BaseController
	validator       validation.IValidator
	login           *login.Handler
	userRepository  user.IUserRepository
	tokenRepository token.ITokenRepository
}

func NewAuthController(
	base *controller.BaseController,
	login *login.Handler,
	validator validation.IValidator,
	userRepository user.IUserRepository,
	tokenRepository token.ITokenRepository,
) *Controller {
	return &Controller{
		BaseController:  base,
		validator:       validator,
		login:           login,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}

// Login godoc
// @Summary Login endpoint
// @tags Login
// @description Login endpoint
// @Accept json
// @Produce json
// @Param dto.LoginRequest body dto.LoginRequest true "Dto для логина"
// @Success      200  {object}  dto.DataLoginResponse
// @Failure      default  {object}  error.HTTPError "The body of any response with an error"
// @Router /api/login [post]
func (authController *Controller) Login(ctx fiber.Ctx) error {
	var loginRequest *dto.LoginRequest

	if err := ctx.Bind().Body(&loginRequest); err != nil {
		return eris.New(err.Error())
	}

	if err := authController.validator.Validate(loginRequest); err != nil {
		return err
	}

	authToken, err := authController.login.Execute(
		ctx,
		(*login.Command)(loginRequest),
	)
	if err != nil {
		return err
	}

	parse, err := time.Parse("2006-01-02 15:04:05", authToken.ExpiresAt)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    authToken.Token,
		Expires:  parse,
		HTTPOnly: true,     // Защита от XSS
		Secure:   true,     // Только HTTPS в production
		SameSite: "Strict", // Защита от CSRF
		Path:     "/",      // Доступно для всех путей
	})

	return ctx.Status(fiber.StatusOK).JSON(&dto.DataLoginResponse{
		Token:     authToken.Token,
		TokenType: authToken.TokenType,
		ExpiresAt: authToken.ExpiresAt,
	})
}

// Logout godoc
// @Summary Logout endpoint
// @tags Logout
// @description Logout endpoint
// @Accept json
// @Produce json
// @Security BearerToken
// @Success      200  {string}  string "message: Successfully logged out."
// @Failure      default  {object}  error.HTTPError "The body of any response with an error"
// @Router /api/logout [get]
func (authController *Controller) Logout(ctx fiber.Ctx) error {
	userCtx := ctx.Locals("user").(*jwt.Token)
	claimsObj := userCtx.Claims.(jwt.MapClaims)

	if err := authController.tokenRepository.RevokeByClaims(
		ctx,
		claimsObj); err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
	})

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message":  "Успешный выход из системы.",
		"redirect": "/login",
	})
}
