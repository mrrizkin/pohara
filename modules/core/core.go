package core

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/core/cache"
	"github.com/mrrizkin/pohara/modules/core/database"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/scheduler"
	"github.com/mrrizkin/pohara/modules/core/server"
	"github.com/mrrizkin/pohara/modules/core/session"
	"github.com/mrrizkin/pohara/modules/core/validator"
)

var Module = fx.Module("core",
	fx.Provide(
		cache.NewRessetto,
		logger.NewZeroLog,
		database.NewGormDB,
		scheduler.NewScheduler,
		validator.NewValidator,
		server.NewServer,
		session.NewSession,
	),

	fx.Decorate(scheduler.LoadSchedules, server.SetupRouter),
	fx.Invoke(scheduler.StartScheduler, server.StartServer),
)
