package job

import (
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/mrrizkin/pohara/internal/scheduler"
	"go.uber.org/fx"
)

type UserScheduler struct {
	log ports.Logger
}

type SchedulerDependencies struct {
	fx.In

	Log ports.Logger
}

type SchedulerResult struct {
	fx.Out

	UserScheduler *UserScheduler
}

func Scheduler(deps SchedulerDependencies) SchedulerResult {
	return SchedulerResult{
		UserScheduler: &UserScheduler{
			log: deps.Log.Scope("user_scheduler"),
		},
	}
}

func (us *UserScheduler) Schedule(s *scheduler.Scheduler) {
	s.Add("* * * * *", func() {
		us.log.Info("I ran every minutes")
	})
}
