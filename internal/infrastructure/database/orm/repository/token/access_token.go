package token

import (
	common "BaseProjectGolang/internal/constant"
	"BaseProjectGolang/internal/domain/tokendmn"
	"BaseProjectGolang/internal/domain/userdmn"
	userModel "BaseProjectGolang/internal/infrastructure/database/orm/model/userModel"
	"BaseProjectGolang/internal/infrastructure/database/query"
	"BaseProjectGolang/pkg/converter"
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/soner3/flora"
)

type AccessTokenRepository struct {
	flora.Component
}

func NewAccessTokenRepository() *AccessTokenRepository {
	return &AccessTokenRepository{}
}

func (tokenRepository *AccessTokenRepository) Create(ctx context.Context, token *tokendmn.AccessToken) error {
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
		Model(&tokendmn.AccessToken{}).
		Where("id = ?", claims["Jti"].(string)).
		Update("Revoked", true).
		Error
}

func (tokenRepository *AccessTokenRepository) DeleteByUser(ctx context.Context, userObj *userdmn.User) error {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	return qb.Current.
		Model(&tokendmn.AccessToken{}).
		Where("user_id = ?", userObj.ID).
		Unscoped().
		Delete(&tokendmn.AccessToken{}).
		Error
}
