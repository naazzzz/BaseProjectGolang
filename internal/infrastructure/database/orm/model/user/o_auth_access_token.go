package user

import (
	"BaseProjectGolang/internal/infrastructure/database/query"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"time"

	common "BaseProjectGolang/internal/common/constant"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
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

func (token *OAuthAccessToken) GenerateID(ctx context.Context) (*OAuthAccessToken, error) {
	qb := ctx.Value(common.QbCtxKey).(*query.Builder)
	n := 40

	var b []byte

	var existToken *OAuthAccessToken

	for ok := true; ok; ok = existToken.ID != "" {
		b = make([]byte, n)
		if _, err := rand.Read(b); err != nil {
			return nil, eris.New(err.Error())
		}

		token.ID = fmt.Sprintf("%X", b)

		result := qb.Original.Model(existToken).Where("id = ?", token.ID).Find(&existToken)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, eris.Wrap(&fiber.Error{Code: fiber.StatusNotFound, Message: result.Error.Error()}, "")
		}
	}

	return token, nil
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
