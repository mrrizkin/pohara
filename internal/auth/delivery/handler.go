package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"go.uber.org/fx"
)

type AuthHandler struct {
	inertia *inertia.Inertia
}

type HandlerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
}

type HandlerResult struct {
	fx.Out

	Auth *AuthHandler
}

func Handler(deps HandlerDependencies) HandlerResult {
	return HandlerResult{
		Auth: &AuthHandler{
			inertia: deps.Inertia,
		},
	}
}

func (h *AuthHandler) LoginPage(ctx *fiber.Ctx) error {
	return h.inertia.Render(ctx, "auth/login")
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusServiceUnavailable)
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusServiceUnavailable)
}
