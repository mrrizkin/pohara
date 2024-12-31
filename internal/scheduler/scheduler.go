// Package scheduler provides a cron-based task scheduling system integrated with fx dependency injection
package scheduler

import (
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

// Schedule defines the interface for tasks that can be scheduled
type Schedule interface {
	Schedule(*Scheduler)
}

// Scheduler wraps cron.Cron with additional functionality
type Scheduler struct {
	*cron.Cron
	logger ports.Logger
}

// Dependencies holds the dependencies for the Scheduler
type Dependencies struct {
	fx.In
	Logger ports.Logger
}

// Result is the fx output container for the Scheduler
type Result struct {
	fx.Out
	Scheduler *Scheduler
}

// Module provides the fx module configuration for the scheduler
var Module = fx.Module("scheduler",
	fx.Provide(NewScheduler),
	fx.Decorate(loadSchedules),
	fx.Invoke(startScheduler),
)

// NewScheduler creates a new Scheduler instance
func NewScheduler(deps Dependencies) Result {
	return Result{
		Scheduler: &Scheduler{
			Cron:   cron.New(),
			logger: deps.Logger,
		},
	}
}

// Add schedules a new task with the given spec
func (s *Scheduler) Add(spec string, cmd func()) error {
	_, err := s.AddFunc(spec, cmd)
	return err
}

// Start begins the scheduler and logs the number of entries
func (s *Scheduler) Start() {
	s.logger.Info("starting scheduler", "entries", len(s.Entries()))
	s.Cron.Start()
}

// SchedulerDependencies contains the scheduler and its registered schedules
type SchedulerDependencies struct {
	fx.In
	Scheduler *Scheduler
	Schedules []Schedule `group:"scheduler"`
}

// loadSchedules registers all provided schedules with the scheduler
func loadSchedules(deps SchedulerDependencies) *Scheduler {
	for _, schedule := range deps.Schedules {
		schedule.Schedule(deps.Scheduler)
	}
	return deps.Scheduler
}

// startScheduler is an fx hook to start the scheduler
func startScheduler(scheduler *Scheduler) {
	scheduler.Start()
}

// AsSchedule annotates a schedule implementation for fx dependency injection
func AsSchedule(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"scheduler"`),
	)
}
