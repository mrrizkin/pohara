package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/app/http/controllers"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/modules/auth"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/core"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/neoweb"
	"github.com/mrrizkin/pohara/routes"
)

func App() *fx.App {
	return fx.New(
		config.Module,
		controllers.Module,
		core.Module,
		neoweb.Module,
		hash.Module,
		repository.Module,
		routes.Module,
		auth.Module,

		fx.WithLogger(func(logger *logger.ZeroLog) fxevent.Logger {
			return logger
		}),
	)
}
