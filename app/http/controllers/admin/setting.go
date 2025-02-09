package admin

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/auth/access"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
)

var (
	ActionSettingProfilePage      = access.NewResourceAction("setting", "profile-page")
	ActionSettingAccountPage      = access.NewResourceAction("setting", "account-page")
	ActionSettingAppearancePage   = access.NewResourceAction("setting", "appearance-page")
	ActionSettingNotificationPage = access.NewResourceAction("setting", "notification-page")
	ActionSettingDisplayPage      = access.NewResourceAction("setting", "display-page")
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
	if !c.auth.Can(ctx, ActionSettingProfilePage, model.MUser{}) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "settings/profile/index")
}

func (c *SettingController) AccountPage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, ActionSettingAccountPage, model.MUser{}) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "settings/account/index")
}

func (c *SettingController) AppearancePage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, ActionSettingAppearancePage, model.MUser{}) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}
	return c.inertia.Render(ctx, "settings/appearance/index")
}

func (c *SettingController) NotificationPage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, ActionSettingNotificationPage, model.MUser{}) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}
	return c.inertia.Render(ctx, "settings/notification/index")
}

func (c *SettingController) DisplayPage(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, ActionSettingDisplayPage, model.MUser{}) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}
	return c.inertia.Render(ctx, "settings/display/index")
}
