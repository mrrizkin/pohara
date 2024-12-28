package user

import (
	"github.com/mrrizkin/pohara/app/user/delivery"
	"github.com/mrrizkin/pohara/app/user/entity"
	"github.com/mrrizkin/pohara/app/user/job"
	"github.com/mrrizkin/pohara/app/user/repository"
	"github.com/mrrizkin/pohara/app/user/usecase"
	"github.com/mrrizkin/pohara/internal/infrastructure/database"
	"github.com/mrrizkin/pohara/internal/scheduler"
	"github.com/mrrizkin/pohara/internal/server"
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Provide(
		repository.Repository,
		usecase.Service,
		job.Scheduler,
		delivery.Handler,

		server.AsApiRouter(delivery.ApiRouter),

		database.AsGormMigration(&entity.User{}),
		scheduler.AsSchedule(func(us *job.UserScheduler) scheduler.Schedule {
			return us
		}),
	),
)
