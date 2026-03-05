package user

import (
	"time"
)

type User struct {
	ID        uint
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Username  string
	Password  string
	Active    bool
}

func (user *User) IsActive() bool {
	return user.Active
}
