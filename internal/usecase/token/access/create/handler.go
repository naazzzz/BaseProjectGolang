package create

import (
	"BaseProjectGolang/internal/config"
	token2 "BaseProjectGolang/internal/domain/token"
	"BaseProjectGolang/internal/domain/user"
	"BaseProjectGolang/internal/usecase/token"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rotisserie/eris"
)

type Handler struct {
	userRepository  user.IUserRepository
	tokenRepository token2.ITokenRepository
	cfg             *config.Config
}

func NewCreateAccessTokenHandler(
	cfg *config.Config,
	userRepository user.IUserRepository,
	tokenRepository token2.ITokenRepository,

) *Handler {
	tokenService := &Handler{
		cfg:             cfg,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}

	return tokenService
}

func (handler *Handler) Execute(
	ctx context.Context,
	domainUser *user.User,
) (*token.Result, error) {
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

	return &token.Result{
		Token:     tokenObj,
		TokenType: "Bearer",
		ExpiresAt: time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
	}, nil
}
