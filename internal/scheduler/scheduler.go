package scheduler

import (
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

type Dependencies struct {
	fx.In

	Log ports.Logger
}

type Result struct {
	fx.Out

	Scheduler *Scheduler
}

func New(deps Dependencies) Result {
	return Result{
		Scheduler: &Scheduler{
			Cron: cron.New(),
			log:  deps.Log,
		},
	}
}

var Module = fx.Module("scheduler",
	fx.Provide(New),
	fx.Decorate(loadScheduler),
	fx.Invoke(func(scheduler *Scheduler) {
		scheduler.Start()
	}),
)

type Scheduler struct {
	*cron.Cron
	log ports.Logger
}

func (s *Scheduler) Add(spec string, cmd func()) {
	s.AddFunc(spec, cmd)
}

func (s *Scheduler) Start() {
	s.log.Info("starting cron", "entries", len(s.Entries()))
	s.Cron.Start()
}

type LoadSchedulerDependencies struct {
	fx.In

	Scheduler *Scheduler
	Schedules []Schedule `group:"scheduler"`
}

func loadScheduler(deps LoadSchedulerDependencies) *Scheduler {
	for _, schedule := range deps.Schedules {
		schedule.Schedule(deps.Scheduler)
	}

	return deps.Scheduler
}
