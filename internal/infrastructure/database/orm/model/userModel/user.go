package userModel

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Username  string `gorm:"uniqueIndex;not null;type:varchar(191)" json:"username"`
	Password  string `gorm:"not null"`
	Active    bool   `json:"active"`

	Tokens []*OAuthAccessToken `gorm:"-" json:"-"`
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
