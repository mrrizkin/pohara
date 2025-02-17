package routes

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/server"
)

var Provide = fx.Provide(
	server.AsWebRouter(WebRouter),
	server.AsApiRouter(ApiRouterV1),
)
