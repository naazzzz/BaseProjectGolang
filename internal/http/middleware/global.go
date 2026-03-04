package middleware

import (
	"fmt"

	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/infrastructure/database"

	"github.com/dromara/carbon/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type GlobalMiddleware struct {
	cfg *config.Config
	db  *database.DataBase
}

func NewGlobalMiddleware(
	cfg *config.Config,
	db *database.DataBase,
) *GlobalMiddleware {
	return &GlobalMiddleware{
		cfg,
		db,
	}
}

func (gl *GlobalMiddleware) SetUpRecover(app *fiber.App) {
	app.Use(
		recover.New(
			recover.Config{
				EnableStackTrace: true,
			},
		),
	)
}

func (gl *GlobalMiddleware) SetUpFiberLogger(app *fiber.App) {
	app.Use(
		logger.New(
			logger.Config{
				TimeFormat: carbon.AtomLayout,
				//Output:     log.Writer(),
				Done: func(_ fiber.Ctx, logString []byte) {
					fmt.Println(string(logString))
				},
			},
		),
	)
}

func (gl *GlobalMiddleware) SetUpStatic(app *fiber.App) {
	app.Get("/assets*", static.New("./web/assets"))

	app.Get("/favicon.ico", func(c fiber.Ctx) error {
		return c.SendFile("assets/favicon.ico")
	})

	app.Get("/static/*", static.New("./web/views/static"))
}
