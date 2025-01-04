package settings

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/settings/delivery"
	"github.com/mrrizkin/pohara/internal/server"
)

var Module = fx.Module("settings",
	fx.Provide(
		delivery.Handler,
		server.AsWebRouter(delivery.WebRouter),
	),
)
