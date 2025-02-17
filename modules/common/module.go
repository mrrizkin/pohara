package common

import (
	"github.com/mrrizkin/pohara/modules/common/hash"
	"go.uber.org/fx"
)

var Module = fx.Module("common",
	fx.Provide(hash.New),
)
