package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/server"
	"github.com/mrrizkin/pohara/internal/web/inertia"
)

func WebRouter(h *AuthHandler, i *inertia.Inertia) server.WebRouter {
	return server.NewWebRouter("/_/auth", func(r fiber.Router) {
		r.Get("/login", i.Middleware(h.LoginPage)).Name("login")
	}, "auth.")
}

func ApiRouter(h *AuthHandler) server.ApiRouter {
	return server.NewApiRouter("v1", "/auth", func(r fiber.Router) {
		r.Get("/login", h.Login).Name("login")
	}, "api.auth.")
}
