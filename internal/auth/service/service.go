package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/internal/auth/access"
	"github.com/mrrizkin/pohara/internal/auth/entity"
)

type Service struct {
	pe *PolicyEnforcement
}

type ServiceDependencies struct {
	PolicyEnforcement *PolicyEnforcement
}

func NewService(deps ServiceDependencies) *Service {
	return &Service{pe: deps.PolicyEnforcement}
}

func (s *Service) IsAllowed(ctx *fiber.Ctx, subject entity.Subject, action access.Action, resource entity.Resource) bool {
	return s.pe.IsAllowed(subject, action, resource, entity.AuthContext{
		IP:       ctx.IP(),
		Time:     time.Now(),
		Location: ctx.Path(),
	})
}
