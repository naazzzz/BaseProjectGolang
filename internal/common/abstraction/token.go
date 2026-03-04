package abstractions

import (
	"context"

	"BaseProjectGolang/internal/infrastructure/database/orm/model/user"

	"github.com/golang-jwt/jwt/v5"
)

type ITokenRepository interface {
	Create(ctx context.Context, token *user.OAuthAccessToken) error
	RevokeByClaims(ctx context.Context, claims jwt.MapClaims) error
	DeleteByUser(ctx context.Context, userObj *user.User) error
}
