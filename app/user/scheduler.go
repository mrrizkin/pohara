package user

import (
	"github.com/mrrizkin/pohara/module/logger"
	"github.com/mrrizkin/pohara/module/scheduler"
	"go.uber.org/fx"
)

type UserScheduler struct {
	log *logger.Logger
}

type SchedulerDependencies struct {
	fx.In

	Log *logger.Logger
}

type SchedulerResult struct {
	fx.Out

	UserScheduler *UserScheduler
}

func Scheduler(deps SchedulerDependencies) SchedulerResult {
	return SchedulerResult{
		UserScheduler: &UserScheduler{
			log: deps.Log,
		},
	}
}

func (us *UserScheduler) Schedule(s *scheduler.Scheduler) {
	s.Add("* * * * *", func() {
		us.log.Info("I ran every minutes")
	})
}
