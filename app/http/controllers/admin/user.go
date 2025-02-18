package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/modules/abac"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/inertia"
	"github.com/mrrizkin/pohara/modules/logger"
)

type UserController struct {
	inertia *inertia.Inertia
	auth    *abac.Authorization
	log     *logger.Logger

	userRepo *repository.UserRepository
}

type UserControllerDependencies struct {
	fx.In

	Inertia       *inertia.Inertia
	Authorization *abac.Authorization
	Logger        *logger.Logger

	UserRepository *repository.UserRepository
}

func NewUserController(deps UserControllerDependencies) *UserController {
	return &UserController{
		inertia:  deps.Inertia,
		auth:     deps.Authorization,
		log:      deps.Logger.Scope("admin_user_controller"),
		userRepo: deps.UserRepository,
	}
}

func (c *UserController) Index(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageUser, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	page := int64(ctx.QueryInt("page", 1))
	limit := int64(ctx.QueryInt("limit", 10))
	searchQ := ctx.Query("q", "")

	search := sql.StringNull()
	if searchQ != "" {
		search.Valid = true
		search.String = searchQ
	}

	result, err := c.userRepo.Find(search, repository.QueryPaginateParams{
		Page:  sql.Int64(page),
		Limit: sql.Int64(limit),
	})
	if err != nil {
		c.log.Error("error get list user", "error", err)
	}

	return c.inertia.Render(ctx, "users/index", gonertia.Props{
		"users": result,
	})
}
