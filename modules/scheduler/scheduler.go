package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/logger"
)

type Scheduler struct {
	core *cron.Cron

	log *logger.Logger
}

type SchedulerDeps struct {
	fx.In

	Logger *logger.Logger
}

// NewScheduler creates a new Scheduler instance
func NewScheduler(lc fx.Lifecycle, deps SchedulerDeps) *Scheduler {
	return &Scheduler{
		core: cron.New(),
		log:  deps.Logger,
	}
}

func startScheduler(lc fx.Lifecycle, scheduler *Scheduler) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			scheduler.start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			scheduler.stop()
			return nil
		},
	})

}

// Add schedules a new task with the given spec
func (s *Scheduler) Add(spec string, cmd func()) error {
	_, err := s.core.AddFunc(spec, cmd)
	return err
}

// Start begins the scheduler and logs the number of entries
func (s *Scheduler) start() {
	s.log.Info("starting scheduler", "entries", len(s.core.Entries()))
	s.core.Start()
}

func (s *Scheduler) stop() {
	s.log.Info("stopping scheduler", "entries", len(s.core.Entries()))
	s.core.Stop()
}
