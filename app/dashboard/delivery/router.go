package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/server"
)

func WebRouter(h *DashboardHandler) server.WebRouter {
	return server.NewWebRouter("/_/", func(r fiber.Router) {
		r.Get("/", h.Index).Name("index")
	}, "dashboard.")
}
