package model

import (
	"github.com/mrrizkin/pohara/modules/auth/access"
	"github.com/mrrizkin/pohara/modules/auth/condition"
)

type CfgPolicy struct {
	ID         uint                `json:"id"         gorm:"primary_key"`
	Name       string              `json:"name"       gorm:"not null,uniqueIndex"`
	Conditions condition.Condition `json:"conditions" gorm:"type:jsonb"`
	Action     access.Action       `json:"action"     gorm:"not null"`
	Effect     access.Effect       `json:"effect"     gorm:"not null"`
	Resource   string              `json:"resource"   gorm:"not null"`
}
