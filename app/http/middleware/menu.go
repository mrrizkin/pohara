package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/modules/auth/access"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
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
		action.PageIntegration,

		action.PageSettingProfile,
		action.PageSettingAccount,
		action.PageSettingAppearance,
		action.PageSettingNotification,
		action.PageSettingDisplay,

		action.PageUser,
	})

	return ctx.Next()
}
