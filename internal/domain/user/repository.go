package user

import (
	"context"
)

type (
	IUserRepository interface {
		Create(ctx context.Context, user *User) (uint, error)
		GetByUsername(ctx context.Context, username string) (*User, error)
		Update(ctx context.Context, user *User) error
		Delete(ctx context.Context, username string, force bool) error
	}
)
