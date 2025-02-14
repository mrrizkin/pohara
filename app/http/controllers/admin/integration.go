package admin

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
)

type IntegrationController struct {
	inertia *inertia.Inertia
	auth    *service.AuthService
}

type IntegrationControllerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
	Auth    *service.AuthService
}

func NewIntegrationController(deps IntegrationControllerDependencies) *IntegrationController {
	return &IntegrationController{
		inertia: deps.Inertia,
		auth:    deps.Auth,
	}
}

func (c *IntegrationController) Index(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageIntegration, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "integrations/index")
}
