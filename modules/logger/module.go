package logger

import (
	"github.com/mrrizkin/pohara/modules/common/config"
	"go.uber.org/fx"
)

var Module = fx.Module("logger",
	config.Load(&Config{}),
	fx.Provide(NewLogger),
)
