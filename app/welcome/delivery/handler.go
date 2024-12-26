package delivery

import (
	"github.com/mrrizkin/pohara/internal/server"
	"go.uber.org/fx"
)

type WelcomeHandler struct {
}

type HandlerResult struct {
	fx.Out

	WelcomeHandler *WelcomeHandler
}

func Handler() HandlerResult {
	return HandlerResult{
		WelcomeHandler: &WelcomeHandler{},
	}
}

func (h *WelcomeHandler) Index(ctx *server.Ctx) error {
	return ctx.Render("welcome", server.Map{})
}
