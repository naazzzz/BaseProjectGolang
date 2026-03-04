package dependency

import (
	"log"

	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/http/controller"
	"BaseProjectGolang/internal/http/controller/auth"
	"BaseProjectGolang/internal/http/middleware"
	authMiddleware "BaseProjectGolang/internal/http/middleware/auth"
	"BaseProjectGolang/internal/http/middleware/context"
	"BaseProjectGolang/internal/infrastructure/database"
	logUtil "BaseProjectGolang/pkg/log"

	"github.com/valyala/fasthttp"
)

type Handlers struct {
	BaseController *controller.BaseController
	// SwaggerController          *global.SwaggerController
	AuthController *auth.Controller

	GlobalMiddleware *middleware.GlobalMiddleware
	AuthMiddleware   *authMiddleware.JwtAddAuthUserInCtx
	SetupCtxQB       *context.SetupCtxQB
}

func NewHandlers(
	globalMiddleware *middleware.GlobalMiddleware,
	authMiddleware *authMiddleware.JwtAddAuthUserInCtx,
	baseController *controller.BaseController,
	setupCtxQB *context.SetupCtxQB,
	authController *auth.Controller,
) *Handlers {
	return &Handlers{
		BaseController:   baseController,
		GlobalMiddleware: globalMiddleware,
		AuthMiddleware:   authMiddleware,
		SetupCtxQB:       setupCtxQB,
		AuthController:   authController,
	}
}

func InitServicesBeforeDi() (db *database.DataBase, cfg *config.Config) {
	var err error

	cfg, err = config.LoadConfig(false, "")
	if err != nil {
		log.Println(err)
	}

	logger := logUtil.InitLogger(cfg.Logs)

	db, err = database.NewDataBase(cfg, logger)
	if err != nil {
		panic(err)
	}

	return
}

func NewClient() *fasthttp.Client {
	return &fasthttp.Client{}
}
