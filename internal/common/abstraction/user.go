package abstractions

import (
	"context"

	"BaseProjectGolang/internal/http/dto"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/user"
)

type (
	IUserService interface {
		CreateUser(
			ctx context.Context,
			userRequest *dto.UserRequest,
		) (err error)

		UpdateUserByUserRequest(
			ctx context.Context,
			userRequest *dto.UserRequest,
			username string,
		) (err error)

		DeleteUser(
			ctx context.Context,
			username string,
			force bool,
		) (err error)
	}

	IUserRepository interface {
		GetUserByUsername(ctx context.Context, username string) (user *user.User, err error)
		Create(ctx context.Context, user *user.User) (err error)
	}
)
