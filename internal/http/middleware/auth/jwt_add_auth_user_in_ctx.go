package auth

import (
	common "BaseProjectGolang/internal/constant"
	"context"
	"errors"
	"net/http"

	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/infrastructure/database"
	userModel "BaseProjectGolang/internal/infrastructure/database/orm/model/userModel"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type JwtAddAuthUserInCtx struct {
	db  *database.DataBase
	cfg *config.Config
}

func NewJwtAddAuthUserInCtx(db *database.DataBase, cfg *config.Config) *JwtAddAuthUserInCtx {
	return &JwtAddAuthUserInCtx{
		db,
		cfg,
	}
}

func (jwtAddAuthUserInCtx *JwtAddAuthUserInCtx) AddAuthUserInCtx(ctx fiber.Ctx) (err error) {
	userCtx := ctx.Locals("user").(*jwt.Token)
	claimsObj := userCtx.Claims.(jwt.MapClaims)
	authorizedUser := &userModel.User{}
	oAuthAccessToken := &userModel.OAuthAccessToken{}

	if result := jwtAddAuthUserInCtx.db.DatabaseDriver.MustGetGorm().First(oAuthAccessToken, "id = ?", claimsObj["Jti"].(string)); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return eris.Wrap(fiber.NewError(fiber.StatusUnauthorized, http.StatusText(fiber.StatusUnauthorized)), "Token not found")
		}

		return eris.Wrap(result.Error, http.StatusText(fiber.StatusInternalServerError))
	}

	if oAuthAccessToken.Revoked {
		return eris.Wrap(fiber.NewError(fiber.StatusUnauthorized, http.StatusText(fiber.StatusUnauthorized)), "Token revoked")
	}

	// todo хранить в кэше?
	if err = jwtAddAuthUserInCtx.db.DatabaseDriver.MustGetGorm().Model(authorizedUser).Where("id = ?", oAuthAccessToken.UserID).Find(authorizedUser).Error; err != nil {
		return eris.Wrap(fiber.NewError(fiber.StatusUnauthorized, "User not found"), http.StatusText(fiber.StatusUnauthorized))
	}

	ctx.Locals(common.AuthorizedUser, authorizedUser)

	gormCtx := jwtAddAuthUserInCtx.db.DatabaseDriver.MustGetGorm().Statement.Context
	gormCtx = context.WithValue(gormCtx, common.AuthorizedUser, authorizedUser)

	jwtAddAuthUserInCtx.db.DatabaseDriver.MustGetGorm().Statement.Context = gormCtx

	return ctx.Next()
}
