package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/modules/abac/access"
	"github.com/mrrizkin/pohara/modules/inertia"
	"go.uber.org/fx"
)

type MenuMiddleware struct {
	inertia *inertia.Inertia
}

type MenuMiddlewareDependencies struct {
	fx.In

	Inertia *inertia.Inertia
}

func NewMenuMiddleware(deps MenuMiddlewareDependencies) *MenuMiddleware {
	return &MenuMiddleware{
		inertia: deps.Inertia,
	}
}

func (m *MenuMiddleware) AuthorizeMenu(ctx *fiber.Ctx) error {
	m.inertia.AuthorizedMenu(ctx, []access.Action{
		action.PageUser,

		action.PageSettingProfile,
		action.PageSettingAccount,
		action.PageSettingAppearance,
		action.PageSettingNotification,
		action.PageSettingDisplay,

		action.PageSystemSettingBranding,
		action.PageSystemSettingIntegration,
		action.PageSystemSettingGeneral,
		action.PageSystemSettingSecurity,
		action.PageSystemSettingLocalization,
		action.PageSystemSettingNotification,
		action.PageSystemSettingBackup,
		action.PageSystemSettingPurge,

		action.PageSystemAuthRole,
		action.PageSystemAuthPolicy,
	})

	return ctx.Next()
}
