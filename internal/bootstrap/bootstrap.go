package bootstrap

import (
	"BaseProjectGolang/internal/command"
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/dependency"
	errorHandler "BaseProjectGolang/internal/http/error"
	"BaseProjectGolang/internal/http/route"
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/internal/infrastructure/database/orm/plugin"
	"BaseProjectGolang/pkg/fasthttpmock"
	"BaseProjectGolang/pkg/log"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/soner3/flora"
)

// App represents the application.
type (
	App struct {
		flora.Component
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

// NewApp creates a new App instance.
func NewApp(
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

// setUpSettings initializes the settings for the application.
func (app *App) setUpSettings() (err error) {
	app.setUpDefaultInstance()
	app.setUpGlobalMiddleware()

	if err = app.RegistrationAllGormPlugins(); err != nil {
		return
	}

	app.setUpScheduler()

	return
}

// RegistrationAllGormPlugins registers all GORM plugins with the application.
func (app *App) RegistrationAllGormPlugins() (err error) {
	if err = app.Plugin.RegisterPluginsInGorm(); err != nil {
		return
	}

	return
}

// setUpScheduler sets up the scheduler for the application.
func (app *App) setUpScheduler() {
	go app.Scheduler.Schedule(context.Background())
}

const timeoutCtxValue = 5 * time.Second

// Stop stops the application gracefully.
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

// setUpDefaultInstance initializes the default instance of Fiber.
func (app *App) setUpDefaultInstance() {
	app.FiberInstance = fiber.New(fiber.Config{
		ErrorHandler: errorHandler.NewErrorHandler(),
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		AppName:      app.Config.AppWorkMode,
	})
}

// getTemplatePath returns the path to the templates directory
func getTemplatePath() string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	return filepath.Join(wd, "web", "views", "templates")
}

// setUpGlobalMiddleware sets up global middleware for Fiber instance
func (app *App) setUpGlobalMiddleware() {
	app.Handlers.GlobalMiddleware.SetUpStatic(app.FiberInstance)
	app.Handlers.GlobalMiddleware.SetUpFiberLogger(app.FiberInstance)
	app.Handlers.GlobalMiddleware.SetUpRecover(app.FiberInstance)
}

// runDefaultHTTPServer starts the HTTP server.
func (app *App) runDefaultHTTPServer() (err error) {
	return app.FiberInstance.Listen(":" + app.Config.AppPort)
}

// Run starts the application
func (app *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	// Start HTTP server in a separate goroutine.
	go func() {
		errCh <- app.runDefaultHTTPServer()
	}()

	// Add another goroutine to handle graceful shutdown.

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		app.Logger.Info("Shutting down servers")
		return nil
	}
}
