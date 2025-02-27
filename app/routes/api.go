package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/http/controllers/api"
	"github.com/mrrizkin/pohara/app/service"
	"github.com/mrrizkin/pohara/modules/server"
)

type ApiRouterV1Dependencies struct {
	fx.In

	AuthService *service.AuthService

	User *api.UserController
}

func ApiRouterV1(deps ApiRouterV1Dependencies) server.ApiRouter {
	return server.NewApiRouter("v1", func(r fiber.Router) {
		authenticated := r.Group("/", deps.AuthService.Authenticated)

		user := authenticated.Group("user").Name("user.")
		user.Get("/", deps.User.UserFind).Name("index")
		user.Get("/:id", deps.User.UserFindByID).Name("show")
		user.Post("/", deps.User.UserCreate).Name("create")
		user.Put("/:id", deps.User.UserUpdate).Name("update")
		user.Delete("/:id", deps.User.UserDelete).Name("destroy")

		r.Get("/*", func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusNotFound)
		}).Name("error.not-found")
	}, "api.v1.")
}
