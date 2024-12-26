package welcome

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/module/server"
	"github.com/mrrizkin/pohara/module/template"
)

var Module = fx.Module("user",
	fx.Provide(
		server.AsWebRouter(func(t *template.Template) server.WebRouter {
			return server.WebRouter{
				Prefix: "/",
				Router: func(r fiber.Router) {
					r.Get("/", func(ctx *fiber.Ctx) error {
						return t.Render(ctx, "welcome", fiber.Map{})
					}).Name("welcome")
				},
			}
		}),
	),
)
