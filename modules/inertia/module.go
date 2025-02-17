package inertia

import "go.uber.org/fx"

var Module = fx.Module("inertia",
	fx.Provide(New),
)
