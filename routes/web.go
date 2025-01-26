package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/http/controllers"
	"github.com/mrrizkin/pohara/app/http/controllers/admin"
	"github.com/mrrizkin/pohara/app/http/middleware"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
	"github.com/mrrizkin/pohara/modules/server"
)

type WebRouterDependencies struct {
	fx.In

	Inertia *inertia.Inertia

	AuthMiddleware *middleware.AuthMiddleware

	Dashboard *admin.DashboardController
	Setting   *admin.SettingController

	Setup *controllers.SetupController

	Auth    *controllers.AuthController
	Welcome *controllers.WelcomeController
}

func WebRouter(deps WebRouterDependencies) server.WebRouter {
	return server.NewWebRouter("/", func(r fiber.Router) {
		r.Get("/", deps.Welcome.Index).Name("welcome")

		webAdmin := r.Group("/_/", deps.Inertia.Middleware)
		webAdmin.Get("/", deps.AuthMiddleware.Authenticated, deps.Dashboard.Index).
			Name("dashboard.index")

		setup := webAdmin.Group("/setup").Name("setup.")
		setup.Get("/", deps.Setup.Index).Name("index")
		setup.Post("/", deps.Setup.Setup).Name("setup")

		auth := webAdmin.Group("/auth").Name("auth.")
		auth.Get("/login", deps.Auth.LoginPage).Name("login")
		auth.Post("/login", deps.Auth.Login).Name("login")
		auth.Get("/register", deps.Auth.RegisterPage).Name("register")
		auth.Post("/register", deps.Auth.Register).Name("register")
		auth.Post("/logout", deps.AuthMiddleware.Authenticated, deps.Auth.Logout).Name("logout")

		setting := webAdmin.Group("/settings", deps.AuthMiddleware.Authenticated).Name("setting.")
		setting.Get("/", deps.Setting.ProfilePage).Name("index")
		setting.Get("/account", deps.Setting.AccountPage).Name("account")
		setting.Get("/appearance", deps.Setting.AppearancePage).Name("appearance")
		setting.Get("/notifications", deps.Setting.NotificationsPage).Name("notifications")
		setting.Get("/display", deps.Setting.DisplayPage).Name("display")

		webAdmin.Get("*", func(ctx *fiber.Ctx) error {
			return deps.Inertia.Render(ctx, "error/not-found", gonertia.Props{})
		}).Name("error.not-found")
	}, "web.")
}
