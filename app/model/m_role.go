package model

import "github.com/mrrizkin/pohara/modules/common/sql"

// Role represents a set of permissions
type MRole struct {
	ID          uint             `json:"id"          gorm:"primaryKey"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CreatedAt   sql.TimeNullable `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt   sql.TimeNullable `json:"updated_at"  gorm:"autoUpdateTime"`
	Policies    []CfgPolicy      `json:"-"           gorm:"-:all"`
}

func (MRole) TableName() string {
	return "m_role"
}
