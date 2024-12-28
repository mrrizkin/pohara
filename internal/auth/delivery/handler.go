package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/server"
	"go.uber.org/fx"
)

type AuthHandler struct {
}

type HandlerResult struct {
	fx.Out

	Auth *AuthHandler
}

func Handler() HandlerResult {
	return HandlerResult{
		Auth: &AuthHandler{},
	}
}

func (h *AuthHandler) LoginPage(ctx *server.Ctx) error {
	return ctx.InertiaRender("auth/login")
}

func (h *AuthHandler) Login(ctx *server.Ctx) error {
	return ctx.SendStatus(fiber.StatusServiceUnavailable)
}

func (h *AuthHandler) Logout(ctx *server.Ctx) error {
	return ctx.SendStatus(fiber.StatusServiceUnavailable)
}
