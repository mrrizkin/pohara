package model

import (
	"github.com/mrrizkin/pohara/modules/abac/access"
	"github.com/mrrizkin/pohara/modules/common/sql"
)

type CfgPolicy struct {
	ID        uint               `json:"id"         gorm:"primaryKey"`
	Name      string             `json:"name"`
	Condition sql.StringNullable `json:"condition"`
	Action    access.Action      `json:"action"`
	Effect    access.Effect      `json:"effect"`
	Resource  sql.StringNullable `json:"resource"`
	CreatedAt sql.TimeNullable   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt sql.TimeNullable   `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CfgPolicy) TableName() string {
	return "cfg_policy"
}

func (p CfgPolicy) GetCondition() sql.StringNullable {
	return p.Condition
}
func (p CfgPolicy) GetAction() access.Action {
	return p.Action
}
func (p CfgPolicy) GetEffect() access.Effect {
	return p.Effect
}
func (p CfgPolicy) GetResource() sql.StringNullable {
	return p.Resource
}
