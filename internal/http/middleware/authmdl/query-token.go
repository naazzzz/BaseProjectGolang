package authmdl

import (
	"BaseProjectGolang/internal/config"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
)

type QueryToken struct {
	cfg *config.Config
}

func NewQueryToken(cfg *config.Config) *QueryToken {
	return &QueryToken{
		cfg: cfg,
	}
}

func (queryToken *QueryToken) AuthorizeViaHeader(c fiber.Ctx) error {
	authKey := c.Get("Authorization")
	if authKey == "" {
		return eris.Wrap(fiber.NewError(http.StatusUnauthorized, ""), "Unauthorized")
	}

	return c.Next()
}
