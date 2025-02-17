package vite

import "go.uber.org/fx"

var Module = fx.Module("vite",
	fx.Provide(New),
)
