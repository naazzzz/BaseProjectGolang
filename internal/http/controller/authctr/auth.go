package authctr

import (
	"BaseProjectGolang/internal/domain/tokendmn"
	"BaseProjectGolang/internal/domain/userdmn"
	"BaseProjectGolang/internal/http/controller"
	"BaseProjectGolang/internal/http/dto"
	"BaseProjectGolang/internal/usecase/userusc/auth/login"
	"BaseProjectGolang/internal/validation"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rotisserie/eris"
	"github.com/soner3/flora"
)

type Controller struct {
	flora.Component `flora:"constructor=NewAuthController"`
	*controller.BaseController
	validator       validation.IValidator
	login           *login.LoginHandler
	userRepository  userdmn.IUserRepository
	tokenRepository tokendmn.ITokenRepository
}

func NewAuthController(
	base *controller.BaseController,
	login *login.LoginHandler,
	validator validation.IValidator,
	userRepository userdmn.IUserRepository,
	tokenRepository tokendmn.ITokenRepository,
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
	userCtx := ctx.Locals("userdmn").(*jwt.Token)
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
