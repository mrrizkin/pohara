package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/web/template"
	"go.uber.org/fx"
)

type WelcomeHandler struct {
	view *template.Template
}

type HandlerDependencies struct {
	fx.In

	Template *template.Template
}

type HandlerResult struct {
	fx.Out

	WelcomeHandler *WelcomeHandler
}

func Handler(deps HandlerDependencies) HandlerResult {
	return HandlerResult{
		WelcomeHandler: &WelcomeHandler{
			view: deps.Template,
		},
	}
}

func (h *WelcomeHandler) Index(ctx *fiber.Ctx) error {
	return h.view.Render(ctx, "welcome", fiber.Map{})
}
