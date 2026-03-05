package token

import (
	common "BaseProjectGolang/internal/constant"
	"BaseProjectGolang/internal/domain/token"
	"BaseProjectGolang/internal/domain/user"
	userModel "BaseProjectGolang/internal/infrastructure/database/orm/model/userModel"
	"BaseProjectGolang/internal/infrastructure/database/query"
	"BaseProjectGolang/pkg/converter"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenRepository struct {
}

func NewTokenRepository() *AccessTokenRepository {
	return &AccessTokenRepository{}
}

func (tokenRepository *AccessTokenRepository) Create(ctx context.Context, token *token.AccessToken) error {
	entity, err := converter.TypeConverter[*userModel.OAuthAccessToken](token)
	if err != nil {
		return err
	}

	qb := ctx.Value(common.QbCtxKey).(*query.Builder)
	return qb.Current.
		Create(entity).
		Error
}

func (tokenRepository *AccessTokenRepository) RevokeByClaims(ctx context.Context, claims jwt.MapClaims) error {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	return qb.Current.
		Model(&token.AccessToken{}).
		Where("id = ?", claims["Jti"].(string)).
		Update("Revoked", true).
		Error
}

func (tokenRepository *AccessTokenRepository) DeleteByUser(ctx context.Context, userObj *user.User) error {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	return qb.Current.
		Model(&token.AccessToken{}).
		Where("user_id = ?", userObj.ID).
		Unscoped().
		Delete(&token.AccessToken{}).
		Error
}
