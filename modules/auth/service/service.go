package service

import (
	"errors"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/auth/access"
	"github.com/mrrizkin/pohara/modules/auth/repository"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/session"
)

type AuthService struct {
	sessionStore *session.Store
	log          *logger.ZeroLog
	authRepo     *repository.AuthRepository
}

type AuthServiceDependencies struct {
	fx.In

	SessionStore *session.Store
	Logger       *logger.ZeroLog
	AuthRepo     *repository.AuthRepository
}

func NewAuthService(deps AuthServiceDependencies) *AuthService {
	return &AuthService{
		sessionStore: deps.SessionStore,
		log:          deps.Logger,
		authRepo:     deps.AuthRepo,
	}
}

func (a *AuthService) Authenticated(ctx *fiber.Ctx) error {
	sess, err := a.sessionStore.Get(ctx)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	authUID := sess.Get("__auth-uid")
	if authUID == nil {
		return fiber.ErrUnauthorized
	}

	uid, ok := authUID.(uint)
	if !ok {
		return fiber.ErrUnauthorized
	}

	userContext, err := a.authRepo.GetUserContext(ctx.Context(), uid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.log.Error("failed to get user attributes", "error", err)
		return fiber.ErrUnauthorized
	}

	ctx.Locals("__auth-user-context", userContext)
	return nil
}

func (a *AuthService) Login(ctx *fiber.Ctx, uid uint) error {
	sess, err := a.sessionStore.Get(ctx)
	if err != nil {
		return err
	}

	sess.Set("__auth-uid", uid)
	err = sess.Save()
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Logout(ctx *fiber.Ctx) error {
	sess, err := a.sessionStore.Get(ctx)
	if err != nil {
		return err
	}

	sess.Delete("__auth-uid")
	err = sess.Save()
	if err != nil {
		return err
	}

	return nil
}

type Resource interface {
	TableName() string
}

func (a *AuthService) Can(
	ctx *fiber.Ctx,
	action access.Action,
	resource Resource,
) bool {
	userContext, ok := ctx.Locals("__auth-user-context").(*repository.UserContext)
	if !ok {
		return false
	}

	if userContext == nil {
		return false
	}

	env := map[string]interface{}{
		"User":           userContext.User,
		"User.Attribute": userContext.UserAttribute,
		"User.Setting":   userContext.UserSetting,
		"Resource":       resource,
	}

	for _, role := range userContext.Roles {
		for _, policy := range role.Policies {
			// special super user check
			if policy.Action == access.ActionGeneralAll && policy.Resource == "all" {
				return true
			}

			if policy.Resource == resource.TableName() && policy.Action == action {
				if !policy.Condition.Valid {
					return policy.Effect == "allow"
				}

				condition := strings.TrimSpace(policy.Condition.String)
				if len(condition) == 0 {
					return policy.Effect == "allow"
				}

				env["Role"] = role.Name
				result, err := a.evaluatePolicy(policy, env)
				if err != nil {
					a.log.Error("failed to evaluate policy", "error", err)
					return false
				}
				if result {
					return policy.Effect == "allow"
				}
			}
		}
	}

	return false
}

func (a *AuthService) evaluatePolicy(
	policy model.CfgPolicy,
	env map[string]interface{},
) (bool, error) {
	program, err := expr.Compile(policy.Condition.String, expr.Env(env))
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}

	result, ok := output.(bool)
	if !ok {
		err := errors.New("evaluation output is not boolean")
		a.log.Error("failed to evaluate policy", "error", err)
		return false, err
	}

	return result, nil
}
