package abstractions

import (
	"context"

	"BaseProjectGolang/internal/http/dto"
	"BaseProjectGolang/internal/infrastructure/database/orm/model/user"
)

type IAuthService interface {
	Login(
		ctx context.Context,
		loginRequest *dto.LoginRequest,
	) (*user.User, error)
	CreateTokenForUser(
		ctx context.Context,
		userObj *user.User,
	) (*dto.TokenInfo, error)
}
