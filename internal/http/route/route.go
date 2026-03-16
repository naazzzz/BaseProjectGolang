package route

import (
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/dependency"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	router   *fiber.App
	handlers *dependency.Handlers
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

	// Базовые middleware для всех маршрутов
	baseGroup := fiberInstance.Group("/")
	baseGroup.Use(
		handlers.SetupCtxQB.DefaultQueryBuilderMiddleware,
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
		return c.Redirect().To("/api/example")
	})

	// Публичные API
	publicAPI := baseGroup.Group("/api")
	router.ExampleAPI(publicAPI)
}

// setupProtectedRoutes настраивает защищенные маршруты (с JWT)
func (router *Router) setupProtectedRoutes(baseGroup fiber.Router) {
	protectedGroup := baseGroup.Group("")

	// Защищенные API
	_ = protectedGroup.Group("/api")
}

func (router *Router) ExampleAPI(route fiber.Router) {
	route.Post("/example", router.handlers.ExampleController.Example)
}
