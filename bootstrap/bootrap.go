package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/pohara/app/dashboard"
	"github.com/mrrizkin/pohara/app/errors"
	"github.com/mrrizkin/pohara/app/settings"
	"github.com/mrrizkin/pohara/app/user"
	"github.com/mrrizkin/pohara/app/welcome"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/auth"
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
	return fx.New(
		config.Load(
			&config.App{},
			&config.Database{},
			&config.Session{},
			&config.Inertia{},
		),

		fx.Provide(
			hashing.New,
			session.New,
			validator.New,
		),

		errors.Module,
		auth.Module,
		settings.Module,
		cache.Module,
		dashboard.Module,
		database.Module,
		inertia.Module,
		logger.Module,
		scheduler.Module,
		server.Module,
		template.Module,
		user.Module,
		vite.Module,
		welcome.Module,

		fx.WithLogger(func(log ports.Logger) fxevent.Logger {
			return log
		}),
	)
}
