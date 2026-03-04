package user

import (
	"BaseProjectGolang/internal/infrastructure/database/query"
	"context"

	common "BaseProjectGolang/internal/common/constant"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/user"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenRepository struct {
}

func NewTokenRepository() *AccessTokenRepository {
	return &AccessTokenRepository{}
}

func (tokenRepository *AccessTokenRepository) Create(ctx context.Context, token *user.OAuthAccessToken) error {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	return qb.Current.
		Create(*token).
		Error
}

func (tokenRepository *AccessTokenRepository) RevokeByClaims(ctx context.Context, claims jwt.MapClaims) error {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	return qb.Current.
		Model(&user.OAuthAccessToken{}).
		Where("id = ?", claims["Jti"].(string)).
		Update("Revoked", true).
		Error
}

func (tokenRepository *AccessTokenRepository) DeleteByUser(ctx context.Context, userObj *user.User) error {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	return qb.Current.
		Model(&user.OAuthAccessToken{}).
		Where("user_id = ?", userObj.ID).
		Unscoped().
		Delete(&user.OAuthAccessToken{}).
		Error
}
