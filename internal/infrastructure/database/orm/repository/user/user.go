package user

import (
	"BaseProjectGolang/internal/infrastructure/database/query"
	"context"
	"net/http"

	common "BaseProjectGolang/internal/common/constant"
	errorHandler "BaseProjectGolang/internal/http/error"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/user"

	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
)

type Repository struct {
}

func NewUserRepository() *Repository {
	return &Repository{}
}

func (userRepository *Repository) GetUserByUsername(ctx context.Context, username string) (user *user.User, err error) {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)
	if err = qb.Original.First(&user, "username = ?", username).Error; err != nil {
		return nil, eris.Wrap(fiber.NewError(fiber.StatusForbidden, err.Error()), "Wrong username or password")
	}

	return
}

func (userRepository *Repository) Create(ctx context.Context, user *user.User) (err error) {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)
	if tx := qb.Current.Create(user); tx.Error != nil {
		return eris.Wrap(err, "Database error")
	} else if tx.RowsAffected == 0 {
		return errorHandler.NewHTTPError(
			http.StatusText(fiber.StatusConflict),
			fiber.StatusConflict,
			err.Error(),
			nil,
			nil,
		)
	}

	return
}
