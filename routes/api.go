package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/app/http/controllers/api"
	"github.com/mrrizkin/pohara/modules/core/server"
	"go.uber.org/fx"
)

type ApiRouterV1Dependencies struct {
	fx.In

	User *api.UserController
}

func ApiRouterV1(deps ApiRouterV1Dependencies) server.ApiRouter {
	return server.NewApiRouter("v1", func(r fiber.Router) {

		user := r.Group("user").Name("user.")
		user.Get("/", deps.User.UserFind).Name("index")
		user.Get("/:id", deps.User.UserFindByID).Name("show")
		user.Post("/", deps.User.UserCreate).Name("create")
		user.Put("/:id", deps.User.UserUpdate).Name("update")
		user.Delete("/:id", deps.User.UserDelete).Name("destroy")

	}, "api.v1.")
}
