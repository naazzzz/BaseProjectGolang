package user

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	User *User
}

func (token *OAuthAccessToken) GetClaimsObj() *jwt.MapClaims {
	return &jwt.MapClaims{
		"Aud":    strconv.FormatUint(uint64(token.ClientID), 10),
		"Jti":    token.ID,
		"Iat":    token.CreatedAt.Unix(),
		"Nbf":    token.UpdatedAt.Unix(),
		"Exp":    token.ExpiresAt.Unix(),
		"Sub":    strconv.FormatUint(uint64(token.UserID), 10),
		"Scopes": token.Scopes,
	}
}
