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

	Tokens []*OAuthAccessToken
}

func (user *User) IsActive() bool {
	return user.Active
}

func (user *User) NewToken() *OAuthAccessToken {
	now := time.Now()
	yearAfter := time.Now().AddDate(1, 0, 0)

	return &OAuthAccessToken{
		ID:        "",
		UserID:    user.ID,
		ClientID:  DefaultMockClientIDNumber,
		Name:      "PersonalAccessToken",
		Scopes:    "[]",
		Revoked:   false,
		CreatedAt: &now,
		UpdatedAt: &now,
		ExpiresAt: &yearAfter,
	}
}
