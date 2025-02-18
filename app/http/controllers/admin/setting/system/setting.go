package system

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/modules/abac/service"
	"github.com/mrrizkin/pohara/modules/inertia"
)

type SystemSettingController struct {
	inertia *inertia.Inertia
	auth    *service.Authorization
}

type SystemSettingControllerDependencies struct {
	fx.In

	Inertia *inertia.Inertia
	Auth    *service.Authorization
}

func NewSystemSettingController(deps SystemSettingControllerDependencies) *SystemSettingController {
	return &SystemSettingController{
		inertia: deps.Inertia,
		auth:    deps.Auth,
	}
}

func (c *SystemSettingController) Branding(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemSettingBranding, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/setting/branding/index")
}

func (c *SystemSettingController) Integration(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemSettingIntegration, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/setting/integrations/index")
}

func (c *SystemSettingController) General(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemSettingGeneral, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/setting/general/index")
}

func (c *SystemSettingController) Security(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemSettingSecurity, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/setting/security/index")
}

func (c *SystemSettingController) Localization(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemSettingLocalization, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/setting/localization/index")
}

func (c *SystemSettingController) Notification(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemSettingNotification, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/setting/notification/index")
}

func (c *SystemSettingController) Backup(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemSettingBackup, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/setting/backu/index")
}

func (c *SystemSettingController) Purge(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemSettingPurge, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/setting/purge/index")
}
