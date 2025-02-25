package admin

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/app/response"
	"github.com/mrrizkin/pohara/modules/abac"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/inertia"
	"github.com/mrrizkin/pohara/modules/logger"
	"github.com/mrrizkin/pohara/modules/server/utils"
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
		log:      deps.Logger.System().Scope("admin_user_controller"),
		userRepo: deps.UserRepository,
	}
}

func (c *UserController) Index(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageUser, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "users/index")
}

type RequestQueryUserDatatable struct {
	Page  sql.Int64Nullable `query:"page,default=1"`
	Limit sql.Int64Nullable `query:"limit,default=10"`

	Name sql.StringNullable `query:"name"`

	Sort sql.StringNullable `query:"sort"`
}

func (c *UserController) Datatable(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageUser, nil) {
		return ctx.Status(fiber.StatusForbidden).
			JSON(response.ErrorMsg("forbidden", "you don't have access to this resource"))
	}

	filter := new(RequestQueryUserDatatable)

	if err := utils.ParseQueryParams(ctx, filter); err != nil {
		return err
	}

	result, err := c.userRepo.Find(filter.Name, repository.QueryPaginateParams{
		Page:  filter.Page,
		Limit: filter.Limit,
	})
	if err != nil {
		c.log.Error("error get list user", "error", err)
	}

	return ctx.JSON(response.Success("success retrieve datatable", result))
}
