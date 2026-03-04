package auth

import (
	"context"
	"errors"
	"time"

	_interface "BaseProjectGolang/internal/common/abstraction"
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/http/dto"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/user"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rotisserie/eris"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepository  _interface.IUserRepository
	tokenRepository _interface.ITokenRepository
	cfg             *config.Config
}

func NewAuthService(
	cfg *config.Config,
	userRepository _interface.IUserRepository,
	tokenRepository _interface.ITokenRepository,

) *Service {
	authService := &Service{
		cfg:             cfg,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}

	return authService
}

func (authService *Service) Login(
	ctx context.Context,
	loginRequest *dto.LoginRequest,
) (*user.User, error) {
	userObj, err := authService.userRepository.GetUserByUsername(ctx, loginRequest.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userObj.Password), []byte(loginRequest.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, eris.Wrap(fiber.ErrForbidden, "Wrong username or password")
		}

		return nil, eris.Wrap(err, "Internal bcrypt error")
	}

	return userObj, nil
}

func (authService *Service) CreateTokenForUser(
	ctx context.Context,
	userObj *user.User,
) (*dto.TokenInfo, error) {
	oAuthAccessToken, err := userObj.NewToken().GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	err = authService.tokenRepository.Create(ctx, oAuthAccessToken)
	if err != nil {
		return nil, err
	}

	tokenObj, err := jwt.NewWithClaims(jwt.SigningMethodHS256, oAuthAccessToken.GetClaimsObj()).SignedString([]byte(authService.cfg.Secure.AuthPrivateKey))
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return &dto.TokenInfo{
		Token:     tokenObj,
		TokenType: "Bearer",
		ExpiresAt: time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
	}, nil
}
