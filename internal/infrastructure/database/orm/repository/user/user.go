package user

import (
	common "BaseProjectGolang/internal/constant"
	domainUser "BaseProjectGolang/internal/domain/user"
	"BaseProjectGolang/internal/infrastructure/database/query"
	"BaseProjectGolang/pkg/converter"
	"context"
	"errors"
	"net/http"
	"time"

	errorHandler "BaseProjectGolang/internal/http/error"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/userModel"

	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type Repository struct {
}

func NewUserRepository() *Repository {
	return &Repository{}
}

func (userRepository *Repository) GetByUsername(ctx context.Context, username string) (result *domainUser.User, err error) {
	var entity *userModel.User

	qb := ctx.Value(common.QbCtxKey).(*query.Builder)
	if err = qb.Original.First(&entity, "username = ?", username).Error; err != nil {
		return nil, eris.Wrap(fiber.NewError(fiber.StatusForbidden, err.Error()), "Wrong username or password")
	}

	return converter.TypeConverter[*domainUser.User](entity)
}

func (userRepository *Repository) Create(ctx context.Context, newUserDomainType *domainUser.User) (createdUserID uint, err error) {
	newUserEntityType, err := converter.TypeConverter[*userModel.User](newUserDomainType)
	if err != nil {
		return 0, err
	}

	qb := ctx.Value(common.QbCtxKey).(*query.Builder)
	if tx := qb.Current.Create(newUserEntityType); tx.Error != nil {
		return 0, eris.Wrap(err, "Database error")
	} else if tx.RowsAffected == 0 {
		return 0, errorHandler.NewHTTPError(
			http.StatusText(fiber.StatusConflict),
			fiber.StatusConflict,
			err.Error(),
			nil,
			nil,
		)
	}

	return newUserEntityType.ID, nil
}

func (userRepository *Repository) Update(
	ctx context.Context,
	user *domainUser.User,
) error {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	entity, err := converter.TypeConverter[*userModel.User](user)
	if err != nil {
		return err
	}

	tx := qb.Current.Model(&userModel.User{}).
		Where("id = ?", entity.ID).
		Updates(entity)

	if tx.Error != nil {
		return eris.Wrap(tx.Error, "update error")
	}

	return nil
}

func (userRepository *Repository) Delete(
	ctx context.Context,
	username string,
	force bool,
) (err error) {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	var userID uint
	if !force {
		// Получаем только ID пользователя, чтобы кинуть событие
		if err = qb.Original.
			Model(&userModel.User{}).
			Select("id").
			Where("username = ?", username).
			First(&userID).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = eris.Wrap(&fiber.Error{Code: fiber.StatusNotFound, Message: err.Error()}, "Record not found")
			} else {
				err = eris.Wrap(err, "Can't get user ID after soft delete")
			}

			return
		}

		if tx := qb.Original.
			Model(&userModel.User{}).
			Where("username = ?", username).
			Update("expiration_date", time.Now()); tx.Error != nil {
			err = eris.Wrap(err, "Can't get user from DBState")
			return
		} else if tx.RowsAffected == 0 {
			err = eris.Wrap(&fiber.Error{Code: fiber.StatusNotFound, Message: err.Error()}, "Record not found")
			return
		}
	} else {
		// get old model for plugin work
		oldUserModel := &userModel.User{}
		if err = qb.Original.
			Model(oldUserModel).
			Where("username = ?", username).
			First(oldUserModel).
			Error; errors.Is(err, gorm.ErrRecordNotFound) {
			err = eris.Wrap(&fiber.Error{Code: fiber.StatusNotFound, Message: err.Error()}, "Record not found")

			return
		} else if err != nil {
			err = eris.Wrap(err, "Can't get user from DBState")
			return
		}

		// получаем ID пользователя перед удалением для ивента
		userID = oldUserModel.ID

		if tx := qb.Original.
			Unscoped().
			Delete(oldUserModel); tx.Error != nil {
			err = eris.Wrap(err, "Can't get user from DBState")
			return
		} else if tx.RowsAffected == 0 {
			err = eris.Wrap(&fiber.Error{Code: fiber.StatusNotFound, Message: err.Error()}, "Record not found")
			return
		}
	}

	return
}
