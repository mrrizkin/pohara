package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/http/controllers"
	"github.com/mrrizkin/pohara/app/http/controllers/admin"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/core/server"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
)

type WebRouterDependencies struct {
	fx.In

	AuthService *service.AuthService
	Inertia     *inertia.Inertia

	Dashboard *admin.DashboardController
	Setting   *admin.SettingController

	Auth    *controllers.AuthController
	Welcome *controllers.WelcomeController
}

func WebRouter(deps WebRouterDependencies) server.WebRouter {
	return server.NewWebRouter("/", func(r fiber.Router) {
		r.Get("/", deps.Welcome.Index).Name("welcome")

		admin := r.Group("/_/", deps.Inertia.Middleware())
		authenticated := admin.Group("/", deps.AuthService.Authenticated)

		auth := admin.Group("/auth").Name("auth.")
		auth.Get("/login", deps.Auth.LoginPage).Name("login")
		auth.Post("/login", deps.Auth.Login).Name("login")
		auth.Get("/register", deps.Auth.RegisterPage).Name("register")
		auth.Post("/register", deps.Auth.Register).Name("register")
		authenticated.Post("/auth/logout", deps.Auth.Logout).Name("auth.logout")

		dashboard := authenticated.Group("/dashboard").Name("dashboard.")
		dashboard.Get("/", deps.Dashboard.Index).Name("index")

		setting := authenticated.Group("/settings").Name("setting.")
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
