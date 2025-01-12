package routes

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/core/server"
)

var Module = fx.Module("routes",
	fx.Provide(
		server.AsWebRouter(WebRouter),
		server.AsApiRouter(ApiRouterV1),
	),
)
