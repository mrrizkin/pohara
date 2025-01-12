package service

import (
	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/auth/access"
)

type PolicyDecisionService struct {
}

func NewPolicyDecisionService() *PolicyDecisionService {
	return &PolicyDecisionService{}
}

func (a *PolicyDecisionService) IsAllowed(
	subject *model.MUser,
	attributes *model.MUserAttribute,
	action access.Action,
	resources interface{},
) bool {
	return false
}
