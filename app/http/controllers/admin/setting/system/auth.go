package system

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/modules/abac/service"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/inertia"
	"github.com/mrrizkin/pohara/modules/logger"
)

type AuthController struct {
	inertia  *inertia.Inertia
	auth     *service.Authorization
	roleRepo *repository.RoleRepository
	log      *logger.Logger
}

type AuthControllerDependencies struct {
	fx.In

	Inertia  *inertia.Inertia
	Auth     *service.Authorization
	RoleRepo *repository.RoleRepository
	Logger   *logger.Logger
}

func NewAuthController(deps AuthControllerDependencies) *AuthController {
	return &AuthController{
		inertia:  deps.Inertia,
		auth:     deps.Auth,
		roleRepo: deps.RoleRepo,
		log:      deps.Logger.Scope("system_setting_auth_controller"),
	}
}

func (c *AuthController) Role(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemAuthRole, nil) {
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

	result, err := c.roleRepo.Find(search, repository.QueryPaginateParams{
		Page:  sql.Int64(page),
		Limit: sql.Int64(limit),
	})
	if err != nil {
		c.log.Error("error get list role", "error", err)
	}

	return c.inertia.Render(ctx, "system/auth/role/index", gonertia.Props{
		"roles": result,
	})
}

func (c *AuthController) Policy(ctx *fiber.Ctx) error {
	if !c.auth.Can(ctx, action.PageSystemAuthPolicy, nil) {
		return c.inertia.Render(ctx.Status(403), "error/forbidden")
	}

	return c.inertia.Render(ctx, "system/auth/policy/index")
}
