package service

import (
	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/auth/access"
	"go.uber.org/fx"
)

type PolicyEnforcementService struct {
	decisionService *PolicyDecisionService
}
type PolicyEnforcementServiceDependencies struct {
	fx.In

	DecisionService *PolicyDecisionService
}

func NewPolicyEnforcementService(
	deps PolicyEnforcementServiceDependencies,
) *PolicyEnforcementService {
	return &PolicyEnforcementService{
		decisionService: deps.DecisionService,
	}
}

func (e *PolicyEnforcementService) IsAllowed(
	subject *model.MUser,
	attributes *model.MUserAttribute,
	action access.Action,
	resources interface{},
) bool {
	return e.decisionService.EvaluateAccess(subject, attributes, action, resources)
}
