package route

import (
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/dependency"
	errorHandler "BaseProjectGolang/internal/http/error"
	"BaseProjectGolang/internal/http/middleware/authmdl"

	"github.com/gofiber/fiber/v3"
	jwtware "github.com/saveblush/gofiber3-contrib/jwt"
)

type Router struct {
	router   *fiber.App
	handlers *dependency.Handlers
	jwtWare  fiber.Handler
}

func NewRoutes(
	cfg *config.Config,
	fiberInstance *fiber.App,
	handlers *dependency.Handlers,
) *Router {
	route := &Router{
		router:   fiberInstance,
		handlers: handlers,
	}

	// Инициализация JWT middleware
	jwtWare := jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(cfg.Secure.AuthPrivateKey)},
		SuccessHandler: handlers.AuthMiddleware.AddAuthUserInCtx,
		ErrorHandler:   errorHandler.NewErrorHandler(),
	})
	route.jwtWare = jwtWare

	// Базовые middleware для всех маршрутов
	baseGroup := fiberInstance.Group("/")
	baseGroup.Use(
		handlers.SetupCtxQB.DefaultQueryBuilderMiddleware,
		authmdl.CookieAuthMiddleware,
	)

	// Настройка маршрутов
	route.setupPublicRoutes(baseGroup)
	route.setupProtectedRoutes(baseGroup)

	return route
}

// setupPublicRoutes настраивает публичные маршруты (без JWT)
func (router *Router) setupPublicRoutes(baseGroup fiber.Router) {
	// Редирект с корня
	baseGroup.Get("/", func(c fiber.Ctx) error {
		return c.Redirect().To("/login")
	})

	// Публичные API
	publicAPI := baseGroup.Group("/api")
	router.LoginAPI(publicAPI)
}

// setupProtectedRoutes настраивает защищенные маршруты (с JWT)
func (router *Router) setupProtectedRoutes(baseGroup fiber.Router) {
	protectedGroup := baseGroup.Group("")
	protectedGroup.Use(router.jwtWare)

	// Защищенные API
	protectedAPI := protectedGroup.Group("/api")
	router.LogoutAPI(protectedAPI)
}

func (router *Router) LoginAPI(route fiber.Router) {
	route.Post("/login", router.handlers.AuthController.Login)
}

func (router *Router) LogoutAPI(route fiber.Router) {
	route.Get("/logout", router.handlers.AuthController.Logout)
}
