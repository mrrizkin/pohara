package entity

import (
	"github.com/mrrizkin/pohara/internal/common/sql"
	"gorm.io/gorm"
)

type Policy struct {
	ID          uint             `json:"id"            gorm:"primary_key"`
	CreatedAt   sql.TimeNullable `json:"created_at"`
	UpdatedAt   sql.TimeNullable `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at"    gorm:"index"`
	Name        string           `json:"name" gorm:"unique;not null;index"`
	Description string           `json:"description"`
	Resource    string           `json:"resource"`
	Action      Action           `json:"action"`
	Effect      Effect           `json:"effect"`
	Priority    uint             `json:"priority"`
	Conditions  Condition        `json:"conditions" gorm:"type:jsonb"`
}

func (Policy) TableName() string {
	return "cfg_policies"
}
