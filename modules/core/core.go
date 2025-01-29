package core

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/core/cache"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/session"
	"github.com/mrrizkin/pohara/modules/core/validator"
	"github.com/mrrizkin/pohara/modules/database"
)

var Module = fx.Module("core",
	fx.Provide(
		cache.NewRessetto,
		logger.NewZeroLog,
		database.NewDatabase,
		validator.NewValidator,
		session.NewSession,
	),
)
