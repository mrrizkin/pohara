package errors

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/errors/delivery"
	"github.com/mrrizkin/pohara/internal/server"
)

var Module = fx.Module("errors",
	fx.Provide(
		delivery.Handler,
		server.AsWebRouter(delivery.WebRouter),
	),
)
