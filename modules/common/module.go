package common

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/common/config"
	"github.com/mrrizkin/pohara/modules/common/hash"
)

var Module = fx.Module("common",
	config.Load(&hash.Config{}),
	fx.Provide(hash.New),
)
