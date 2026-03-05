//go:build wireinject

package app

import (
	"BaseProjectGolang/internal/bootstrap"
	"BaseProjectGolang/internal/command"
	"BaseProjectGolang/internal/config"
	"BaseProjectGolang/internal/dependency"
	token2 "BaseProjectGolang/internal/domain/token"
	user2 "BaseProjectGolang/internal/domain/user"
	"BaseProjectGolang/internal/http/controller"
	"BaseProjectGolang/internal/http/controller/auth"
	"BaseProjectGolang/internal/http/middleware"
	authMiddleware "BaseProjectGolang/internal/http/middleware/auth"
	"BaseProjectGolang/internal/http/middleware/context"
	"BaseProjectGolang/internal/infrastructure/database"
	"BaseProjectGolang/internal/infrastructure/database/orm/plugin"
	"BaseProjectGolang/internal/infrastructure/database/orm/repository/token"
	userRepository "BaseProjectGolang/internal/infrastructure/database/orm/repository/user"
	"BaseProjectGolang/internal/usecase/token/access/create"
	"BaseProjectGolang/internal/usecase/user/auth/login"
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
		wire.Bind(new(validation.IValidator), new(*validation.XValidator)),

		userRepository.NewUserRepository,
		wire.Bind(new(user2.IUserRepository), new(*userRepository.Repository)),
		token.NewTokenRepository,
		wire.Bind(new(token2.ITokenRepository), new(*token.AccessTokenRepository)),

		create.NewCreateAccessTokenHandler,
		login.NewLoginHandler,

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
