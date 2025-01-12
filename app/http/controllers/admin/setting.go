package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
)

type SettingController struct {
	inertia *inertia.Inertia
}

type SettingControllerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
}

func NewSettingController(deps SettingControllerDependencies) *SettingController {
	return &SettingController{
		inertia: deps.Inertia,
	}
}

func (c *SettingController) ProfilePage(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "settings/profile/index", gonertia.Props{})
}

func (c *SettingController) AccountPage(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "settings/account/index", gonertia.Props{})
}

func (c *SettingController) AppearancePage(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "settings/appearance/index", gonertia.Props{})
}

func (c *SettingController) NotificationsPage(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "settings/notifications/index", gonertia.Props{})
}

func (c *SettingController) DisplayPage(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "settings/display/index", gonertia.Props{})
}
