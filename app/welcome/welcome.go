package welcome

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/welcome/delivery"
	"github.com/mrrizkin/pohara/internal/server"
)

var Module = fx.Module("welcome",
	fx.Provide(
		delivery.Handler,
		server.AsWebRouter(delivery.WebRouter),
	),
)
