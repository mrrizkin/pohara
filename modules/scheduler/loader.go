package scheduler

import "go.uber.org/fx"

type Schedule interface {
	Schedule(*Scheduler)
}

// ScheduleDeps contains the scheduler and its registered schedules
type ScheduleDeps struct {
	fx.In
	Scheduler *Scheduler
	Schedules []Schedule `group:"scheduler"`
}

// LoadSchedules registers all provided schedules with the scheduler
func loadSchedules(deps ScheduleDeps) *Scheduler {
	for _, schedule := range deps.Schedules {
		schedule.Schedule(deps.Scheduler)
	}
	return deps.Scheduler
}
