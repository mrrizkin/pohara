package dashboard

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/dashboard/delivery"
	"github.com/mrrizkin/pohara/internal/server"
)

var Module = fx.Module("dashboard",
	fx.Provide(
		delivery.Handler,
		server.AsWebRouter(delivery.WebRouter),
	),
)
