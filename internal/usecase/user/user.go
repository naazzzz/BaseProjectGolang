package user

import (
	_interface "BaseProjectGolang/internal/common/abstraction"
	common "BaseProjectGolang/internal/common/constant"
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/http/dto"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/user"
	"BaseProjectGolang/internal/infrastructure/database/query"
	"BaseProjectGolang/pkg/converter"
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type Service struct {
	cfg            *config.Config
	userRepository _interface.IUserRepository
}

func NewUserService(
	cfg *config.Config,
	userRepository _interface.IUserRepository,
) *Service {
	return &Service{
		cfg:            cfg,
		userRepository: userRepository,
	}
}

// CreateUser Создает нового пользователя
func (userService *Service) CreateUser(
	ctx context.Context,
	userRequest *dto.UserRequest,
) (err error) {
	userObj, err := converter.TypeConverter[*user.User](userRequest)
	if err != nil {
		return eris.Wrap(err, "Type Converter error")
	}

	if err = userService.userRepository.Create(ctx, userObj); err != nil {
		return err
	}

	return
}

// UpdateUserByUserRequest Обновляет все значенния из request.UserRequest DTO для model.User
// так, чтобы отработали все update callback
// todo refactor move to repository
func (userService *Service) UpdateUserByUserRequest(
	ctx context.Context,
	userRequest *dto.UserRequest,
	username string,
) (err error) {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)
	// get old model
	oldUserModel := &user.User{}
	if err = qb.Current.
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

	// convert dto to model.User
	userModel, err := userRequest.ToUserModel()
	if err != nil {
		return err
	}

	// save so that all update callbacks are processed
	if tx := qb.Current.
		Model(oldUserModel).
		Updates(userModel); tx.Error != nil {
		err = eris.Wrap(err, "Can't get user from DBState")
		return
	} else if tx.RowsAffected == 0 {
		err = eris.Wrap(&fiber.Error{Code: fiber.StatusNotFound, Message: err.Error()}, "Record not found")
		return
	}

	if username != userRequest.Username {
		log.Printf("Username was changed for %d old %s >>> new %s", userModel.ID, username, userRequest.Username)
	}

	return
}

// DeleteUser удаляет пользователя по username, флаг force определяет жесткое или мягкое удаление
func (userService *Service) DeleteUser(
	ctx context.Context,
	username string,
	force bool,
) (err error) {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)

	var userID uint
	if !force {
		// Получаем только ID пользователя, чтобы кинуть событие
		if err = qb.Original.
			Model(&user.User{}).
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
			Model(&user.User{}).
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
		oldUserModel := &user.User{}
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
