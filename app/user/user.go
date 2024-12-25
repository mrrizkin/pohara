package user

import (
	"github.com/mrrizkin/pohara/module/scheduler"
	"github.com/mrrizkin/pohara/module/server"
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Provide(
		Repository,
		Service,
		Scheduler,
		Handler,

		server.AsApiRouter(userApiRouter),

		scheduler.AsSchedule(func(us *UserScheduler) scheduler.Schedule {
			return us
		}),
	),
)
