package templ

import "go.uber.org/fx"

var Module = fx.Module("templ",
	fx.Provide(New),
)
