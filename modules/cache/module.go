package cache

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/common/config"
)

var Module = fx.Module("cache",
	config.Load(&Config{}),
	fx.Provide(NewCache),
)
