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

type ClientRouterDependencies struct {
	fx.In

	ClientPage *controllers.ClientPageController
}

type AdminRouterDependencies struct {
	fx.In

	Inertia *inertia.Inertia

	AuthMiddleware *middleware.AuthMiddleware

	Dashboard *admin.DashboardController
	Setting   *admin.SettingController

	Setup *controllers.SetupController

	Auth *controllers.AuthController
}

func ClientRouter(deps ClientRouterDependencies) server.WebRouter {
	return server.NewWebRouter("/", func(r fiber.Router) {
		r.Get("/", deps.ClientPage.HomePage).Name("homepage")
		r.Get("/blog", deps.ClientPage.HomePage).Name("blog")
		r.Get("/pricing", deps.ClientPage.Pricing).Name("pricing")
		r.Get("/contact", deps.ClientPage.Contact).Name("contact")
		r.Get("/faq", deps.ClientPage.Faq).Name("faq")

		r.Get("*", deps.ClientPage.NotFound).Name("error.not-found")
	}, "client.")
}

func AdminRouter(deps AdminRouterDependencies) server.WebRouter {
	return server.NewWebRouter("/_/", func(r fiber.Router) {
		admin := r.Group("/", deps.Inertia.Middleware)
		admin.Get("/", deps.AuthMiddleware.Authenticated, deps.Dashboard.Index).Name("dashboard.index")

		setup := admin.Group("/setup").Name("setup.")
		setup.Get("/", deps.Setup.Index).Name("index")
		setup.Post("/", deps.Setup.Setup).Name("setup")

		auth := admin.Group("/auth").Name("auth.")
		auth.Get("/login", deps.Auth.LoginPage).Name("login")
		auth.Post("/login", deps.Auth.Login).Name("login")
		auth.Get("/register", deps.Auth.RegisterPage).Name("register")
		auth.Post("/register", deps.Auth.Register).Name("register")
		auth.Post("/logout", deps.AuthMiddleware.Authenticated, deps.Auth.Logout).Name("logout")

		setting := admin.Group("/settings", deps.AuthMiddleware.Authenticated).Name("setting.")
		setting.Get("/", deps.Setting.ProfilePage).Name("index")
		setting.Get("/account", deps.Setting.AccountPage).Name("account")
		setting.Get("/appearance", deps.Setting.AppearancePage).Name("appearance")
		setting.Get("/notifications", deps.Setting.NotificationsPage).Name("notifications")
		setting.Get("/display", deps.Setting.DisplayPage).Name("display")

		admin.Get("*", func(ctx *fiber.Ctx) error {
			return deps.Inertia.Render(ctx, "error/not-found", gonertia.Props{})
		}).Name("error.not-found")
	}, "admin.")
}
