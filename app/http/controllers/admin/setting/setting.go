package setting

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/inertia"
)

type SettingController struct {
	inertia *inertia.Inertia
	auth    *service.AuthService
}

type SettingControllerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
	Auth    *service.AuthService
}

func NewSettingController(deps SettingControllerDependencies) *SettingController {
	return &SettingController{
		inertia: deps.Inertia,
		auth:    deps.Auth,
	}
}

func (c *SettingController) ProfilePage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSettingProfile, nil) {
		return c.inertia.Render(ctx.Status(fiber.StatusForbidden), "error/forbidden")
	}

	return c.inertia.Render(ctx, "settings/profile/index")
}

func (c *SettingController) AccountPage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSettingAccount, nil) {
		return c.inertia.Render(ctx.Status(fiber.StatusForbidden), "error/forbidden")
	}

	return c.inertia.Render(ctx, "settings/account/index")
}

func (c *SettingController) AppearancePage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSettingAppearance, nil) {
		return c.inertia.Render(ctx.Status(fiber.StatusForbidden), "error/forbidden")
	}
	return c.inertia.Render(ctx, "settings/appearance/index")
}

func (c *SettingController) NotificationPage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSettingNotification, nil) {
		return c.inertia.Render(ctx.Status(fiber.StatusForbidden), "error/forbidden")
	}
	return c.inertia.Render(ctx, "settings/notification/index")
}

func (c *SettingController) DisplayPage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSettingDisplay, nil) {
		return c.inertia.Render(ctx.Status(fiber.StatusForbidden), "error/forbidden")
	}
	return c.inertia.Render(ctx, "settings/display/index")
}
