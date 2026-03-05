package userModel

import (
	"time"
)

const (
	DefaultMockClientIDNumber = 3
)

type OAuthAccessToken struct {
	ID        string
	UserID    uint
	ClientID  uint
	Name      string
	Scopes    string
	Revoked   bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ExpiresAt *time.Time
	// relations
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (token *OAuthAccessToken) TableName() string {
	return "oauth_access_tokens"
}
