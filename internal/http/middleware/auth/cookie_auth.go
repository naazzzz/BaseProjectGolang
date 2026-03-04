package auth

import "github.com/gofiber/fiber/v3"

func CookieAuthMiddleware(ctx fiber.Ctx) error {
	// Если уже есть Authorization header, пропускаем
	if ctx.Get("Authorization") != "" {
		return ctx.Next()
	}

	// Проверяем куку
	token := ctx.Cookies("auth_token")
	if token != "" {
		// Устанавливаем Authorization header для JWT middleware
		ctx.Request().Header.Set("Authorization", "Bearer "+token)
	}

	return ctx.Next()
}
