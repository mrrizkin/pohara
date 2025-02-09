package model

import (
	"github.com/mrrizkin/pohara/modules/auth/access"
	"github.com/mrrizkin/pohara/modules/common/sql"
)

type CfgPolicy struct {
	ID        uint               `json:"id"         gorm:"primaryKey"`
	Name      string             `json:"name"`
	Condition sql.StringNullable `json:"condition"`
	Action    access.Action      `json:"action"`
	Effect    access.Effect      `json:"effect"`
	Resource  string             `json:"resource"`
	CreatedAt sql.TimeNullable   `json:"created_at"`
	UpdatedAt sql.TimeNullable   `json:"updated_at"`
}

func (CfgPolicy) TableName() string {
	return "cfg_policy"
}
