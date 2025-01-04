package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"
)

type SettingsHandler struct {
	inertia *inertia.Inertia
}

type HandlerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
}

type HandlerResult struct {
	fx.Out

	Settings *SettingsHandler
}

func Handler(deps HandlerDependencies) HandlerResult {
	return HandlerResult{
		Settings: &SettingsHandler{
			inertia: deps.Inertia,
		},
	}
}

func (h *SettingsHandler) ProfilePage(ctx *fiber.Ctx) error {
	return h.inertia.Render(ctx, "settings/profile/index", gonertia.Props{})
}

func (h *SettingsHandler) AccountPage(ctx *fiber.Ctx) error {
	return h.inertia.Render(ctx, "settings/account/index", gonertia.Props{})
}

func (h *SettingsHandler) AppearancePage(ctx *fiber.Ctx) error {
	return h.inertia.Render(ctx, "settings/appearance/index", gonertia.Props{})
}

func (h *SettingsHandler) NotificationsPage(ctx *fiber.Ctx) error {
	return h.inertia.Render(ctx, "settings/notifications/index", gonertia.Props{})
}

func (h *SettingsHandler) DisplayPage(ctx *fiber.Ctx) error {
	return h.inertia.Render(ctx, "settings/display/index", gonertia.Props{})
}
