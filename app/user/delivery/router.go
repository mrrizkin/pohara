package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/server"
)

func ApiRouter(h *UserHandler) server.ApiRouter {
	return server.NewApiRouter("v1", "/user", func(r fiber.Router) {
		r.Get("/", h.UserFind).Name("index")
		r.Get("/:id", h.UserFindByID).Name("show")
		r.Post("/", h.UserCreate).Name("create")
		r.Put("/:id", h.UserUpdate).Name("update")
		r.Delete("/:id", h.UserDelete).Name("destroy")
	}, "api.user.")
}
