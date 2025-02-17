package cli

import "go.uber.org/fx"

var Module = fx.Module("cli",
	fx.Provide(New),
	fx.Decorate(loadCommand),
	fx.Invoke(execute),
)
