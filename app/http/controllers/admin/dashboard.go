package admin

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/inertia"
)

type DashboardController struct {
	inertia *inertia.Inertia
}

type DashboardControllerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
}

func NewDashboardController(deps DashboardControllerDependencies) *DashboardController {
	return &DashboardController{
		inertia: deps.Inertia,
	}
}

func (c *DashboardController) Index(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "dashboard/index")
}
