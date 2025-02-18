package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/pohara/app"
	"github.com/mrrizkin/pohara/app/service"
	"github.com/mrrizkin/pohara/database/migration"
	"github.com/mrrizkin/pohara/modules/abac"
	"github.com/mrrizkin/pohara/modules/cache"
	"github.com/mrrizkin/pohara/modules/cli"
	"github.com/mrrizkin/pohara/modules/common"
	"github.com/mrrizkin/pohara/modules/database"
	"github.com/mrrizkin/pohara/modules/inertia"
	"github.com/mrrizkin/pohara/modules/logger"
	"github.com/mrrizkin/pohara/modules/scheduler"
	"github.com/mrrizkin/pohara/modules/server"
	"github.com/mrrizkin/pohara/modules/session"
	"github.com/mrrizkin/pohara/modules/templ"
	"github.com/mrrizkin/pohara/modules/validator"
	"github.com/mrrizkin/pohara/modules/vite"
)

func App() *fx.App {
	return fx.New(
		app.Module,
		abac.Module,
		cache.Module,
		cli.Module,
		common.Module,
		database.Module,
		inertia.Module,
		logger.Module,
		migration.Module,
		scheduler.Module,
		server.Module,
		session.Module,
		templ.Module,
		validator.Module,
		vite.Module,

		fx.Provide(func(authService *service.AuthService) abac.AuthService {
			return authService
		}),

		fx.WithLogger(func(logger *logger.Logger) fxevent.Logger {
			return logger
		}),
	)
}
