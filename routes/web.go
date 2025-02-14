package routes

import (
	"github.com/gofiber/fiber/v2"
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
	MenuMiddleware *middleware.MenuMiddleware

	Dashboard   *admin.DashboardController
	Setting     *admin.SettingController
	User        *admin.UserController
	Integration *admin.IntegrationController

	Setup *controllers.SetupController

	Auth *controllers.AuthController

	ClientPage *controllers.ClientPageController
}

func WebRouter(deps WebRouterDependencies) server.WebRouter {
	return server.NewWebRouter("/", func(r fiber.Router) {
		r.Get("/", deps.ClientPage.HomePage).Name("homepage")
		r.Get("/blog", deps.ClientPage.HomePage).Name("blog")
		r.Get("/pricing", deps.ClientPage.Pricing).Name("pricing")
		r.Get("/contact", deps.ClientPage.Contact).Name("contact")
		r.Get("/faq", deps.ClientPage.Faq).Name("faq")

		auth := r.Group("/_/auth", deps.Inertia.Middleware).Name("auth.")
		auth.Get("/login", deps.Auth.LoginPage).Name("login")
		auth.Post("/login", deps.Auth.Login).Name("login")
		auth.Get("/register", deps.Auth.RegisterPage).Name("register")
		auth.Post("/register", deps.Auth.Register).Name("register")
		auth.Post("/logout", deps.AuthMiddleware.Authenticated, deps.Auth.Logout).Name("logout")

		setup := r.Group("/_/setup", deps.Inertia.Middleware).Name("setup.")
		setup.Get("/", deps.Setup.Index).Name("index")
		setup.Post("/", deps.Setup.Setup).Name("setup")

		authenticated := r.Group(
			"/_/",
			deps.Inertia.Middleware,
			deps.AuthMiddleware.Authenticated,
			deps.MenuMiddleware.AuthorizeMenu,
		)
		authenticated.Get("/", deps.Dashboard.Index).Name("dashboard.index")

		user := authenticated.Group("/users").Name("user.")
		user.Get("/", deps.User.Index).Name("index")

		integration := authenticated.Group("/integrations").Name("integration.")
		integration.Get("/", deps.Integration.Index).Name("index")

		setting := authenticated.Group("/settings").Name("setting.")
		setting.Get("/", deps.Setting.ProfilePage).Name("index")
		setting.Get("/account", deps.Setting.AccountPage).Name("account")
		setting.Get("/appearance", deps.Setting.AppearancePage).Name("appearance")
		setting.Get("/notifications", deps.Setting.NotificationPage).Name("notifications")
		setting.Get("/display", deps.Setting.DisplayPage).Name("display")

		r.Get("/_/*", func(ctx *fiber.Ctx) error {
			return deps.Inertia.Render(ctx.Status(fiber.StatusNotFound), "error/not-found")
		}).Name("error.not-found")

		r.Get("*", deps.ClientPage.NotFound).Name("error.not-found")
	}, "admin.")
}
