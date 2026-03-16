package login

import (
	"BaseProjectGolang/internal/config"
	token2 "BaseProjectGolang/internal/domain/tokendmn"
	"BaseProjectGolang/internal/domain/userdmn"
	"BaseProjectGolang/internal/usecase/tokenusc"
	"BaseProjectGolang/internal/usecase/tokenusc/access/create"
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
	"github.com/soner3/flora"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	flora.Component
	userRepository    userdmn.IUserRepository
	tokenRepository   token2.ITokenRepository
	createAccessToken *create.CreateAccessTokenHandler
	cfg               *config.Config
}

func NewLoginHandler(
	cfg *config.Config,
	userRepository userdmn.IUserRepository,
	tokenRepository token2.ITokenRepository,
	createAccessToken *create.CreateAccessTokenHandler,
) *LoginHandler {
	authService := &LoginHandler{
		cfg:               cfg,
		userRepository:    userRepository,
		tokenRepository:   tokenRepository,
		createAccessToken: createAccessToken,
	}

	return authService
}

func (handler *LoginHandler) Execute(
	ctx context.Context,
	cmd *Command,
) (tokenInfo *tokenusc.Result, err error) {
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
