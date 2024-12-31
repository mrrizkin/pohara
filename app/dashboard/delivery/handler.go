package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"
)

type DashboardHandler struct {
	inertia *inertia.Inertia
}

type HandlerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
}

type HandlerResult struct {
	fx.Out

	Dashboard *DashboardHandler
}

func Handler(deps HandlerDependencies) HandlerResult {
	return HandlerResult{
		Dashboard: &DashboardHandler{
			inertia: deps.Inertia,
		},
	}
}

func (h *DashboardHandler) Index(ctx *fiber.Ctx) error {
	return h.inertia.Render(ctx, "dashboard/index", gonertia.Props{
		"text": "Pohara the most battery included starterkit",
	})
}
