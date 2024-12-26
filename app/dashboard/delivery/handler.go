package delivery

import (
	"github.com/mrrizkin/pohara/internal/server"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"
)

type DashboardHandler struct {
}

type HandlerResult struct {
	fx.Out

	Dashboard *DashboardHandler
}

func Handler() HandlerResult {
	return HandlerResult{
		Dashboard: &DashboardHandler{},
	}
}

func (h *DashboardHandler) Index(ctx *server.Ctx) error {
	return ctx.InertiaRender("dashboard/index", gonertia.Props{
		"text": "Pohara the most battery included starterkit",
	})
}
