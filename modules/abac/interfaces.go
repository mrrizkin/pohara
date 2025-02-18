package abac

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/abac/access"
	"github.com/mrrizkin/pohara/modules/common/sql"
)

type AuthService interface {
	GetActionSpecialAll() access.Action
	GetSubjectName() string
	GetSubject(ctx *fiber.Ctx) (Subject, error)
	GetEnv(ctx *fiber.Ctx) (any, error)
}

type Subject interface {
	GetPolicies() ([]Policy, error)
}

type Policy interface {
	GetCondition() sql.StringNullable
	GetAction() access.Action
	GetEffect() access.Effect
	GetResource() sql.StringNullable
}

type Resource interface {
	ResourceName() sql.StringNullable
}
