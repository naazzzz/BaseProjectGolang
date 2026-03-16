package dependency

import (
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/http/controller"
	"BaseProjectGolang/internal/http/controller/authctr"
	"BaseProjectGolang/internal/http/middleware"
	authMiddleware "BaseProjectGolang/internal/http/middleware/authmdl"
	"BaseProjectGolang/internal/http/middleware/context"
	log2 "BaseProjectGolang/pkg/log"
	"log"

	"github.com/soner3/flora"
	"github.com/valyala/fasthttp"
)

type Handlers struct {
	flora.Component
	BaseController *controller.BaseController
	// SwaggerController          *global.SwaggerController
	AuthController *authctr.Controller

	GlobalMiddleware *middleware.GlobalMiddleware
	AuthMiddleware   *authMiddleware.JwtAddAuthUserInCtx
	SetupCtxQB       *context.SetupCtxQB
}

func NewHandlers(
	globalMiddleware *middleware.GlobalMiddleware,
	authMiddleware *authMiddleware.JwtAddAuthUserInCtx,
	baseController *controller.BaseController,
	setupCtxQB *context.SetupCtxQB,
	authController *authctr.Controller,
) *Handlers {
	return &Handlers{
		BaseController:   baseController,
		GlobalMiddleware: globalMiddleware,
		AuthMiddleware:   authMiddleware,
		SetupCtxQB:       setupCtxQB,
		AuthController:   authController,
	}
}

type Start struct {
	flora.Configuration
}

func (*Start) NewStart() (cfg *config.Config) {
	var err error

	cfg, err = config.NewConfig(false, "")
	if err != nil {
		log.Println(err)
	}

	return
}

type ClientConfig struct {
	flora.Configuration
}

//flora:primary
func (c *ClientConfig) ProvideClient() *fasthttp.Client {
	return &fasthttp.Client{}
}

type LoggerConfig struct {
	flora.Configuration
}

//flora:primary
func (loggerCfg *LoggerConfig) ProvideLoggerConfig(cfg *config.Config) log2.LoggerConfig {
	return cfg.Logs
}
