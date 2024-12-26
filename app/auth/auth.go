package auth

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/auth/delivery"
	"github.com/mrrizkin/pohara/internal/server"
)

var Module = fx.Module("auth",
	fx.Provide(
		delivery.Handler,
		server.AsWebRouter(delivery.WebRouter),
	),
)
