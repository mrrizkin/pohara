package abac

import (
	"errors"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/abac/access"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/logger"
)

type Authorization struct {
	log     *logger.Logger
	service AuthService
}

type AuthorizationServiceDeps struct {
	fx.In

	Logger  *logger.Logger
	Service AuthService
}

func NewAuthorizationService(deps AuthorizationServiceDeps) *Authorization {
	return &Authorization{
		log:     deps.Logger,
		service: deps.Service,
	}
}

func (a *Authorization) Can(
	ctx *fiber.Ctx,
	action access.Action,
	resource Resource,
) bool {
	if a.service == nil {
		return false
	}

	subjectName := a.service.GetSubjectName()
	subject, err := a.service.GetSubject(ctx)
	if err != nil {
		return false
	}

	if subject == nil {
		return false
	}

	policies, err := subject.GetPolicies()
	if err != nil {
		return false
	}

	if policies == nil {
		return false
	}

	env, err := a.service.GetEnv(ctx)
	if err != nil {
		return false
	}

	data := map[string]interface{}{
		subjectName: subject,
		"Env":       env,
	}

	resourceName := sql.StringNull()
	if resource != nil {
		resourceName = resource.ResourceName()
		if resourceName.Valid {
			data["Resource"] = resource
		}
	}

	for _, policy := range policies {
		// special super user check
		if policy.GetAction() == a.service.GetActionSpecialAll() {
			return true
		}

		if policy.GetResource() == resourceName && policy.GetAction() == action {
			condition := policy.GetCondition()
			if !condition.Valid || len(strings.TrimSpace(condition.String)) == 0 {
				return policy.GetEffect() == access.EffectAllow
			}

			result, err := a.evaluatePolicy(policy, data)
			if err != nil {
				a.log.Error("failed to evaluate policy", "error", err)
				return false
			}
			if result {
				return policy.GetEffect() == access.EffectAllow
			}
		}
	}

	return false
}

func (a *Authorization) evaluatePolicy(
	policy Policy,
	data map[string]interface{},
) (bool, error) {
	// we can optimize this using precompile method.. and mapped it with policy.id as a key
	program, err := expr.Compile(policy.GetCondition().String, expr.Env(data))
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, data)
	if err != nil {
		return false, err
	}

	result, ok := output.(bool)
	if !ok {
		return false, errors.New("evaluation output is not boolean")
	}

	return result, nil
}
