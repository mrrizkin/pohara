package templ

import (
	"context"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/neoweb/vite"
)

type templContext string

var (
	FiberContextKey templContext = "fiber-context"
)

type Templ struct {
	config *config.Config
	vite   *vite.Vite
}

type Dependencies struct {
	fx.In

	Config *config.Config
	Vite   *vite.Vite
}

func New(deps Dependencies) *Templ {
	return &Templ{
		config: deps.Config,
		vite:   deps.Vite,
	}
}

func (t *Templ) Render(ctx *fiber.Ctx, component templ.Component) error {
	ctx.Set("Content-Type", "text/html")
	ctx.Locals("config", t.config)
	ctx.Locals("vite", t.vite)
	c := context.WithValue(context.Background(), FiberContextKey, ctx)
	return component.Render(c, ctx.Response().BodyWriter())
}
