package welcome

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/module/server"
)

var Module = fx.Module("welcome",
	fx.Provide(
		server.AsWebRouter(func() server.WebRouter {
			return server.NewWebRouter("/", func(r *server.Router) {
				r.Get("/", func(ctx *server.Ctx) error {
					return ctx.Render("welcome", fiber.Map{})
				}).Name("index")
			}, "welcome")
		}),
	),
)
