package session

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/common/config"
	sessionConfig "github.com/mrrizkin/pohara/modules/session/config"
)

var Module = fx.Module("session",
	config.Load(&sessionConfig.Config{}),
	fx.Provide(NewSession),
)
