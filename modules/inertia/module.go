package inertia

import (
	"github.com/mrrizkin/pohara/modules/common/config"
	"go.uber.org/fx"
)

var Module = fx.Module("inertia",
	config.Load(&Config{}),
	fx.Provide(New),
)
