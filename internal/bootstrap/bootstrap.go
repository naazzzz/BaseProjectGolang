package bootstrap

import (
	"BaseProjectGolang/internal/command"
	"BaseProjectGolang/internal/infrastructure/database/orm/plugin"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/dependency"
	errorHandler "BaseProjectGolang/internal/http/error"
	"BaseProjectGolang/internal/http/route"
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/pkg/fasthttpmock"
	"BaseProjectGolang/pkg/log"

	"github.com/gofiber/fiber/v3"
)

type (
	App struct {
		FiberInstance *fiber.App
		Route         *route.Router
		Config        *config.Config
		Handlers      *dependency.Handlers
		Database      *database.DataBase
		Scheduler     *command.Scheduler
		Client        *fasthttpmock.WrapClient
		Logger        *log.Logger
		Plugin        plugin.IRegistrationPlugin
	}
)

func Bootstrap(
	conf *config.Config,
	handlers *dependency.Handlers,
	database *database.DataBase,
	scheduler *command.Scheduler,
	client *fasthttpmock.WrapClient,
	logger *log.Logger,
	plugin plugin.IRegistrationPlugin,
) (app *App, err error) {
	app = &App{
		Config:    conf,
		Handlers:  handlers,
		Database:  database,
		Scheduler: scheduler,
		Client:    client,
		Logger:    logger,
		Plugin:    plugin,
	}

	if err = app.setUpSettings(); err != nil {
		return
	}

	app.Route = route.NewRoutes(
		conf,
		app.FiberInstance,
		handlers,
	)

	return
}

func (app *App) setUpSettings() (err error) {
	app.setUpDefaultInstance()
	app.setUpGlobalMiddleware()

	if err = app.RegistrationAllGormPlugins(); err != nil {
		return
	}

	app.setUpScheduler()

	return
}

func (app *App) RegistrationAllGormPlugins() (err error) {
	if err = app.Plugin.RegisterPluginsInGorm(); err != nil {
		return
	}

	return
}

func (app *App) setUpScheduler() {
	go app.Scheduler.Schedule(context.Background())
}

const timeoutCtxValue = 5 * time.Second

func (app *App) Stop() {
	if app.Scheduler != nil {
		app.Scheduler.Stop() // Остановите шедулер
	}

	if app.FiberInstance != nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutCtxValue)
		defer cancel()

		if err := app.FiberInstance.ShutdownWithContext(ctx); err != nil {
			app.Logger.Logger.Printf("Error shutting down Fiber app: %v", err)
		} // Остановите сервер
	}

	if app.Logger != nil {
		err := app.Logger.CloseLogger()
		if err != nil {
			return
		}
	}
}

func (app *App) setUpDefaultInstance() {
	app.FiberInstance = fiber.New(fiber.Config{
		ErrorHandler: errorHandler.NewErrorHandler(),
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		AppName:      app.Config.AppWorkMode,
	})
}

func getTemplatePath() string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	return filepath.Join(wd, "web", "views", "templates")
}

func (app *App) setUpGlobalMiddleware() {
	app.Handlers.GlobalMiddleware.SetUpStatic(app.FiberInstance)
	app.Handlers.GlobalMiddleware.SetUpFiberLogger(app.FiberInstance)
	app.Handlers.GlobalMiddleware.SetUpRecover(app.FiberInstance)
}

func (app *App) Run() (err error) {
	return app.FiberInstance.Listen(":" + app.Config.AppPort)
}
