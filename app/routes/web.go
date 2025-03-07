package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/http/controllers"
	"github.com/mrrizkin/pohara/app/http/controllers/admin"
	"github.com/mrrizkin/pohara/app/http/controllers/admin/setting"
	"github.com/mrrizkin/pohara/app/http/controllers/admin/setting/system"
	"github.com/mrrizkin/pohara/app/http/middleware"
	"github.com/mrrizkin/pohara/modules/inertia"
	"github.com/mrrizkin/pohara/modules/server"
)

type WebRouterDependencies struct {
	fx.In

	Inertia *inertia.Inertia

	AuthMiddleware *middleware.AuthMiddleware
	MenuMiddleware *middleware.MenuMiddleware

	Dashboard *admin.DashboardController
	User      *admin.UserController

	Setting           *setting.SettingController
	SystemSetting     *system.SystemSettingController
	SettingSystemAuth *system.AuthController

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
		user.Get("/datatable", deps.User.Datatable).Name("datatable")

		setting := authenticated.Group("/settings").Name("setting.")
		setting.Get("/", deps.Setting.ProfilePage).Name("index")
		setting.Get("/account", deps.Setting.AccountPage).Name("account")
		setting.Get("/appearance", deps.Setting.AppearancePage).Name("appearance")
		setting.Get("/notifications", deps.Setting.NotificationPage).Name("notifications")
		setting.Get("/display", deps.Setting.DisplayPage).Name("display")

		system := authenticated.Group("/system").Name("system.")
		systemSetting := system.Group("/setting").Name("setting.")
		systemSetting.Get("/branding", deps.SystemSetting.Branding).Name("branding")
		systemSetting.Get("/integrations", deps.SystemSetting.Integration).Name("integration")
		systemSetting.Get("/general", deps.SystemSetting.General).Name("general")
		systemSetting.Get("/security", deps.SystemSetting.Security).Name("security")
		systemSetting.Get("/localization", deps.SystemSetting.Localization).Name("localization")
		systemSetting.Get("/notification", deps.SystemSetting.Notification).Name("notification")
		systemSetting.Get("/backup", deps.SystemSetting.Backup).Name("backup")
		systemSetting.Get("/purge", deps.SystemSetting.Purge).Name("purge")

		authSetting := system.Group("/auth").Name("auth.")
		authSetting.Get("/role", deps.SettingSystemAuth.Role).Name("role")
		authSetting.Get("/role/datatable", deps.SettingSystemAuth.RoleDatatable).Name("role.datatable")
		authSetting.Get("/role/:id/policies", deps.SettingSystemAuth.RolePolicies).Name("role.policies")
		authSetting.Get("/policy", deps.SettingSystemAuth.Policy).Name("policy")
		authSetting.Get("/policy/datatable", deps.SettingSystemAuth.PolicyDatatable).Name("policy.datatable")
		authSetting.Get("/policy/:id/roles", deps.SettingSystemAuth.PolicyRoles).Name("policy.roles")

		r.Get("/_/*", func(ctx *fiber.Ctx) error {
			if inertia.IsInertiaRequest(ctx) {
				return deps.Inertia.Render(ctx.Status(fiber.StatusNotFound), "error/not-found")
			}

			return ctx.SendStatus(fiber.StatusNotFound)
		}).Name("error.not-found")

		r.Get("*", deps.ClientPage.NotFound).Name("error.not-found")
	}, "admin.")
}
