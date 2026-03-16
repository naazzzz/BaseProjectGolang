package tokendmn

import (
	"BaseProjectGolang/internal/domain/userdmn"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type ITokenRepository interface {
	Create(ctx context.Context, token *AccessToken) error
	RevokeByClaims(ctx context.Context, claims jwt.MapClaims) error
	DeleteByUser(ctx context.Context, userObj *userdmn.User) error
}
