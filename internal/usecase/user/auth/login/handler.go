package login

import (
	"BaseProjectGolang/internal/config"
	token2 "BaseProjectGolang/internal/domain/token"
	"BaseProjectGolang/internal/domain/user"
	"BaseProjectGolang/internal/usecase/token"
	"BaseProjectGolang/internal/usecase/token/access/create"
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userRepository    user.IUserRepository
	tokenRepository   token2.ITokenRepository
	createAccessToken *create.Handler
	cfg               *config.Config
}

func NewLoginHandler(
	cfg *config.Config,
	userRepository user.IUserRepository,
	tokenRepository token2.ITokenRepository,
	createAccessToken *create.Handler,
) *Handler {
	authService := &Handler{
		cfg:               cfg,
		userRepository:    userRepository,
		tokenRepository:   tokenRepository,
		createAccessToken: createAccessToken,
	}

	return authService
}

func (handler *Handler) Execute(
	ctx context.Context,
	cmd *Command,
) (tokenInfo *token.Result, err error) {
	userObj, err := handler.userRepository.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userObj.Password), []byte(cmd.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, eris.Wrap(fiber.ErrForbidden, "Wrong username or password")
		}

		return nil, eris.Wrap(err, "Internal bcrypt error")
	}

	tokenInfo, err = handler.createAccessToken.Execute(ctx, userObj)
	if err != nil {
		return nil, err
	}

	return
}
