package scheduler

import "go.uber.org/fx"

type Schedule interface {
	Schedule(*Scheduler)
}

func AsSchedule(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"scheduler"`),
	)
}
