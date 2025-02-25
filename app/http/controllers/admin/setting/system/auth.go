package system

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

type AuthController struct {
	inertia *inertia.Inertia
	auth    *abac.Authorization
	log     *logger.Logger

	roleRepo *repository.RoleRepository
}

type AuthControllerDependencies struct {
	fx.In

	Inertia       *inertia.Inertia
	Authorization *abac.Authorization
	Logger        *logger.Logger

	RoleRepo *repository.RoleRepository
}

func NewAuthController(deps AuthControllerDependencies) *AuthController {
	return &AuthController{
		inertia:  deps.Inertia,
		auth:     deps.Authorization,
		roleRepo: deps.RoleRepo,
		log:      deps.Logger.System().Scope("system_setting_auth_controller"),
	}
}

func (c *AuthController) Role(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemAuthRole, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/auth/role/index")
}

type RequestQueryRoleDatatable struct {
	Page  sql.Int64Nullable `query:"page,default=1"`
	Limit sql.Int64Nullable `query:"limit,default=10"`

	Name sql.StringNullable `query:"name"`

	Sort sql.StringNullable `query:"sort"`
}

func (c *AuthController) RoleDatatable(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemAuthRole, nil) {
		return ctx.Status(fiber.StatusForbidden).
			JSON(response.ErrorMsg("forbidden", "you don't have access to this resource"))
	}

	filter := new(RequestQueryRoleDatatable)

	if err := utils.ParseQueryParams(ctx, filter); err != nil {
		return err
	}

	result, err := c.roleRepo.Find(filter.Name, repository.QueryPaginateParams{
		Page:  filter.Page,
		Limit: filter.Limit,
	})
	if err != nil {
		c.log.Error("error get list role", "error", err)
	}

	return ctx.JSON(response.Success("success retrieve datatable", result))
}

func (c *AuthController) Policy(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemAuthPolicy, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/auth/policy/index")
}
