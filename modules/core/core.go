package core

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/core/cache"
	"github.com/mrrizkin/pohara/modules/core/database"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/session"
	"github.com/mrrizkin/pohara/modules/core/validator"
)

var Module = fx.Module("core",
	fx.Provide(
		cache.NewRessetto,
		logger.NewZeroLog,
		database.NewGormDB,
		validator.NewValidator,
		session.NewSession,
	),
)
