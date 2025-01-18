package scheduler

import (
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

type Scheduler struct {
	*cron.Cron
	logger *logger.ZeroLog
}

type Dependencies struct {
	fx.In
	Logger *logger.ZeroLog
}

var Module = fx.Module("scheduler",
	fx.Provide(
		NewScheduler,
	),
	fx.Decorate(LoadSchedules),
	fx.Invoke(StartScheduler),
)

// NewScheduler creates a new Scheduler instance
func NewScheduler(deps Dependencies) *Scheduler {
	return &Scheduler{
		Cron:   cron.New(),
		logger: deps.Logger,
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
	Schedules []AdditionalScheduler `group:"scheduler"`
}

// LoadSchedules registers all provided schedules with the scheduler
func LoadSchedules(deps SchedulerDependencies) *Scheduler {
	for _, schedule := range deps.Schedules {
		schedule.scheduler(deps.Scheduler)
	}
	return deps.Scheduler
}

// StartScheduler is an fx hook to start the scheduler
func StartScheduler(scheduler *Scheduler) {
	scheduler.Start()
}

// AsSchedule annotates a schedule implementation for fx dependency injection
func AsSchedule(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"scheduler"`),
	)
}

type AdditionalScheduler struct {
	scheduler func(*Scheduler)
}

func New(fn func(*Scheduler)) AdditionalScheduler {
	return AdditionalScheduler{
		scheduler: fn,
	}
}
