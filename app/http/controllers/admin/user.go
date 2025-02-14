package admin

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
)

type UserController struct {
	inertia *inertia.Inertia
	auth    *service.AuthService
}

type UserControllerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
	Auth    *service.AuthService
}

func NewUserController(deps UserControllerDependencies) *UserController {
	return &UserController{
		inertia: deps.Inertia,
		auth:    deps.Auth,
	}
}

func (c *UserController) Index(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageUser, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "users/index")
}
