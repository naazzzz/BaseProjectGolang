package create

import (
	"BaseProjectGolang/internal/config"
	token2 "BaseProjectGolang/internal/domain/tokendmn"
	"BaseProjectGolang/internal/domain/userdmn"
	"BaseProjectGolang/internal/usecase/tokenusc"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rotisserie/eris"
	"github.com/soner3/flora"
)

type CreateAccessTokenHandler struct {
	flora.Component `flora:"constructor=NewCreateAccessTokenHandler"`
	userRepository  userdmn.IUserRepository
	tokenRepository token2.ITokenRepository
	cfg             *config.Config
}

func NewCreateAccessTokenHandler(
	cfg *config.Config,
	userRepository userdmn.IUserRepository,
	tokenRepository token2.ITokenRepository,

) *CreateAccessTokenHandler {
	tokenService := &CreateAccessTokenHandler{
		cfg:             cfg,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}

	return tokenService
}

func (handler *CreateAccessTokenHandler) Execute(
	ctx context.Context,
	domainUser *userdmn.User,
) (*tokenusc.Result, error) {
	oAuthAccessToken, err := token2.NewAccessToken(domainUser.ID).GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	err = handler.tokenRepository.Create(ctx, oAuthAccessToken)
	if err != nil {
		return nil, err
	}

	tokenObj, err := jwt.NewWithClaims(jwt.SigningMethodHS256, oAuthAccessToken.GetClaimsObj()).SignedString([]byte(handler.cfg.Secure.AuthPrivateKey))
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return &tokenusc.Result{
		Token:     tokenObj,
		TokenType: "Bearer",
		ExpiresAt: time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
	}, nil
}
