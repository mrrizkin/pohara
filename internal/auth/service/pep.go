package service

import (
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/common/helper"
	"github.com/mrrizkin/pohara/internal/ports"
)

type RoleHolder interface {
	Role() entity.Role
}

type Enforcer struct {
	db ports.Database
}

func NewEnforcer(db ports.Database) *Enforcer {
	return &Enforcer{db}
}

func (e *Enforcer) Can(roleHolder RoleHolder, action entity.Action, resource interface{}) bool {
	_, err := helper.StructToMap(resource, "json")
	if err != nil {
		return false
	}

	return false
}
