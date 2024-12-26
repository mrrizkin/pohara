package delivery

import (
	"github.com/mrrizkin/pohara/internal/server"
)

func ApiRouter(h *UserHandler) server.ApiRouter {
	return server.NewApiRouter("v1", "/user", func(r *server.Router) {
		r.Get("/", h.UserFindAll).Name("index")
		r.Get("/:id", h.UserFindByID).Name("show")
		r.Post("/", h.UserCreate).Name("create")
		r.Put("/:id", h.UserUpdate).Name("update")
		r.Delete("/:id", h.UserDelete).Name("destroy")
	}, "api.user.")
}
