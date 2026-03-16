package dependency

import (
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/http/controller"
	"BaseProjectGolang/internal/http/controller/examplectr"
	"BaseProjectGolang/internal/http/middleware"
	"BaseProjectGolang/internal/http/middleware/context"
	log2 "BaseProjectGolang/pkg/log"
	"log"

	"github.com/soner3/flora"
	"github.com/valyala/fasthttp"
)

type Handlers struct {
	flora.Component
	BaseController    *controller.BaseController
	ExampleController *examplectr.ExampleController

	GlobalMiddleware *middleware.GlobalMiddleware
	SetupCtxQB       *context.SetupCtxQB
}

func NewHandlers(
	globalMiddleware *middleware.GlobalMiddleware,
	baseController *controller.BaseController,
	setupCtxQB *context.SetupCtxQB,
	authController *examplectr.ExampleController,
) *Handlers {
	return &Handlers{
		BaseController:    baseController,
		GlobalMiddleware:  globalMiddleware,
		SetupCtxQB:        setupCtxQB,
		ExampleController: authController,
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
