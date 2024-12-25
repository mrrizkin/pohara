package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/module/server"
)

func userApiRouter(h *UserHandler) server.ApiRouter {
	return server.ApiRouter{
		Prefix: "/user",
		Router: func(r fiber.Router) {
			r.Get("/", h.UserFindAll).Name("index")
			r.Get("/:id", h.UserFindByID).Name("show")
			r.Post("/", h.UserCreate).Name("create")
			r.Put("/:id", h.UserUpdate).Name("update")
			r.Delete("/:id", h.UserDelete).Name("destroy")
		},
		Name:    "api.user.",
		Version: "v1",
	}
}
