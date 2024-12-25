package welcome

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/module/inertia"
	"github.com/mrrizkin/pohara/module/server"
)

var Module = fx.Module("user",
	fx.Provide(
		server.AsWebRouter(func(i *inertia.Inertia) server.WebRouter {
			return server.WebRouter{
				Prefix: "/",
				Router: func(r fiber.Router) {
					r.Get("/", func(ctx *fiber.Ctx) error {
						return i.Render(ctx, "welcome/index", gonertia.Props{
							"text": "Pohara Is StarterKit",
						})
					}).Name("welcome")
				},
			}
		}),
	),
)
