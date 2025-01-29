package model

import (
	"github.com/mrrizkin/pohara/modules/auth/access"
	"github.com/mrrizkin/pohara/modules/auth/condition"
)

type CfgPolicy struct {
	ID         uint                `json:"id"         db:"id"`
	Name       string              `json:"name"       db:"name"`
	Conditions condition.Condition `json:"conditions" db:"conditions"`
	Action     access.Action       `json:"action"     db:"action"`
	Effect     access.Effect       `json:"effect"     db:"effect"`
	Resource   string              `json:"resource"   db:"resource"`
}

func (CfgPolicy) TableName() string {
	return "cfg_policy"
}
