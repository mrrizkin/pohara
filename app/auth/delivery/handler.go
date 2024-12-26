package delivery

import (
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

func (h *AuthHandler) Login(ctx *server.Ctx) error {
	return ctx.InertiaRender("auth/login")
}
