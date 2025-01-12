package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/http/controllers"
	"github.com/mrrizkin/pohara/app/http/controllers/admin"
	"github.com/mrrizkin/pohara/modules/core/server"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
)

type WebRouterDependencies struct {
	fx.In

	Inertia *inertia.Inertia

	Dashboard *admin.DashboardController
	Setting   *admin.SettingController

	Welcome *controllers.WelcomeController
}

func WebRouter(deps WebRouterDependencies) server.WebRouter {
	return server.NewWebRouter("/", func(r fiber.Router) {
		r.Get("/", deps.Welcome.Index).Name("welcome")

		dashboard := r.Group("/_/", deps.Inertia.Middleware()).Name("dashboard.")
		dashboard.Get("/", deps.Dashboard.Index).Name("index")

		setting := dashboard.Group("/settings").Name("setting.")
		setting.Get("/", deps.Setting.ProfilePage).Name("index")
		setting.Get("/account", deps.Setting.AccountPage).Name("account")
		setting.Get("/appearance", deps.Setting.AppearancePage).Name("appearance")
		setting.Get("/notifications", deps.Setting.NotificationsPage).Name("notifications")
		setting.Get("/display", deps.Setting.DisplayPage).Name("display")

		r.Get("*", func(ctx *fiber.Ctx) error {
			return deps.Inertia.Render(ctx, "error/not-found", gonertia.Props{})
		}).Name("error.not-found")
	}, "web.")
}
