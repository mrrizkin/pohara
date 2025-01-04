package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"
)

type ErrorsHandler struct {
	inertia *inertia.Inertia
}

type HandlerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
}

type HandlerResult struct {
	fx.Out

	Errors *ErrorsHandler
}

func Handler(deps HandlerDependencies) HandlerResult {
	return HandlerResult{
		Errors: &ErrorsHandler{
			inertia: deps.Inertia,
		},
	}
}

func (h *ErrorsHandler) PageNotFound(ctx *fiber.Ctx) error {
	return h.inertia.Render(ctx, "error/not-found", gonertia.Props{})
}
