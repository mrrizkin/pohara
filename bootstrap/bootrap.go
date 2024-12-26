package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/pohara/app/user"
	"github.com/mrrizkin/pohara/app/welcome"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/models"
	"github.com/mrrizkin/pohara/module/cache"
	"github.com/mrrizkin/pohara/module/database"
	"github.com/mrrizkin/pohara/module/hashing"
	"github.com/mrrizkin/pohara/module/inertia"
	"github.com/mrrizkin/pohara/module/logger"
	"github.com/mrrizkin/pohara/module/scheduler"
	"github.com/mrrizkin/pohara/module/server"
	"github.com/mrrizkin/pohara/module/session"
	"github.com/mrrizkin/pohara/module/template"
	"github.com/mrrizkin/pohara/module/validator"
	"github.com/mrrizkin/pohara/module/vite"
)

func App() *fx.App {
	fx.Decorate()
	return fx.New(
		config.New(),

		fx.Provide(
			logger.New,
			database.New,
			hashing.New,
			inertia.New,
			session.New,
			cache.New,
			validator.New,
		),

		template.Module,
		vite.Module,

		user.Module,
		welcome.Module,

		models.Module,
		scheduler.Module,
		server.Module,

		fx.WithLogger(func(log *logger.Logger) fxevent.Logger {
			return log
		}),
	)
}
