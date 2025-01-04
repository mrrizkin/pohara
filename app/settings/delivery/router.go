package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/server"
	"github.com/mrrizkin/pohara/internal/web/inertia"
)

func WebRouter(h *SettingsHandler, i *inertia.Inertia) server.WebRouter {
	return server.NewWebRouter("/_/settings", func(r fiber.Router) {
		r.Get("/", i.Middleware(h.ProfilePage)).Name("index")
		r.Get("/account", i.Middleware(h.AccountPage)).Name("account")
		r.Get("/appearance", i.Middleware(h.AppearancePage)).Name("appearance")
		r.Get("/notifications", i.Middleware(h.NotificationsPage)).Name("notifications")
		r.Get("/display", i.Middleware(h.DisplayPage)).Name("display")
	}, "settings.profile")
}
