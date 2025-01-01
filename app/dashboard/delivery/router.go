package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/server"
	"github.com/mrrizkin/pohara/internal/web/inertia"
)

func WebRouter(h *DashboardHandler, i *inertia.Inertia) server.WebRouter {
	return server.NewWebRouter("/_/", func(r fiber.Router) {
		r.Get("/", i.Middleware(h.Index)).Name("index")
	}, "dashboard.")
}
