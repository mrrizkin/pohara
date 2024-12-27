package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/pohara/app/auth"
	"github.com/mrrizkin/pohara/app/dashboard"
	"github.com/mrrizkin/pohara/app/user"
	"github.com/mrrizkin/pohara/app/welcome"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/common/hashing"
	"github.com/mrrizkin/pohara/internal/common/validator"
	"github.com/mrrizkin/pohara/internal/infrastructure/cache"
	"github.com/mrrizkin/pohara/internal/infrastructure/database"
	"github.com/mrrizkin/pohara/internal/infrastructure/logger"
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/mrrizkin/pohara/internal/scheduler"
	"github.com/mrrizkin/pohara/internal/server"
	"github.com/mrrizkin/pohara/internal/session"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/mrrizkin/pohara/internal/web/template"
	"github.com/mrrizkin/pohara/internal/web/vite"
)

func App() *fx.App {
	fx.Decorate()
	return fx.New(
		config.New(),

		fx.Provide(
			hashing.New,
			session.New,
			validator.New,
		),

		auth.Module,
		user.Module,
		welcome.Module,
		dashboard.Module,

		logger.Module,
		cache.Module,
		database.Module,
		inertia.Module,
		vite.Module,
		template.Module,

		scheduler.Module,
		server.Module,

		fx.WithLogger(func(log ports.Logger) fxevent.Logger {
			return log
		}),
	)
}
