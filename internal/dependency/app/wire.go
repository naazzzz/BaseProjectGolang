//go:build wireinject

package app

import (
	"BaseProjectGolang/internal/bootstrap"
	"BaseProjectGolang/internal/command"
	_interface "BaseProjectGolang/internal/common/abstraction"
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/dependency"
	"BaseProjectGolang/internal/http/controller"
	"BaseProjectGolang/internal/http/controller/auth"
	"BaseProjectGolang/internal/http/middleware"
	authMiddleware "BaseProjectGolang/internal/http/middleware/auth"
	"BaseProjectGolang/internal/http/middleware/context"
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/internal/infrastructure/database/orm/plugin"
	userRepository "BaseProjectGolang/internal/infrastructure/database/orm/repository/user"
	"BaseProjectGolang/internal/usecase/user"
	authService "BaseProjectGolang/internal/usecase/user/auth"
	"BaseProjectGolang/internal/validation"
	"BaseProjectGolang/pkg/fasthttpmock"
	"BaseProjectGolang/pkg/log"
	"BaseProjectGolang/pkg/mail"

	"github.com/google/wire"
)

func InitializeApp(
	cfg *config.Config,
	db *database.DataBase,
) (*bootstrap.App, error) {
	wire.Build(
		GetLoggerConfig,

		log.InitLogger,

		plugin.NewRegistrationPlugin,
		wire.Bind(new(plugin.IRegistrationPlugin), new(*plugin.RegistrationPlugin)),

		dependency.NewClient,
		fasthttpmock.NewWrapClient,

		validation.NewValidator,
		wire.Bind(new(_interface.IValidator), new(*validation.XValidator)),

		userRepository.NewUserRepository,
		wire.Bind(new(_interface.IUserRepository), new(*userRepository.Repository)),
		userRepository.NewTokenRepository,
		wire.Bind(new(_interface.ITokenRepository), new(*userRepository.AccessTokenRepository)),

		authService.NewAuthService,
		wire.Bind(new(_interface.IAuthService), new(*authService.Service)),
		user.NewUserService,
		wire.Bind(new(_interface.IUserService), new(*user.Service)),

		controller.NewBaseController,

		auth.NewAuthController,

		command.NewScheduler,

		middleware.NewGlobalMiddleware,
		authMiddleware.NewJwtAddAuthUserInCtx,
		context.NewSetupCtxQB,

		dependency.NewHandlers,

		bootstrap.Bootstrap,
	)

	return &bootstrap.App{}, nil
}

func GetMailConfig(cfg *config.Config) *mail.GomailServiceConfig {
	return cfg.Mail
}

func GetLoggerConfig(cfg *config.Config) log.LoggerConfig {
	return cfg.Logs
}
