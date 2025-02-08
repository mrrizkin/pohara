package routes

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/server"
)

var Module = fx.Module("routes",
	fx.Provide(
		server.AsWebRouter(ClientRouter),
		server.AsWebRouter(AdminRouter),
		server.AsApiRouter(ApiRouterV1),
	),
)
