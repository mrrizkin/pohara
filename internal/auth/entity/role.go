package entity

import (
	"github.com/mrrizkin/pohara/internal/common/sql"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint             `json:"id"          gorm:"primary_key"`
	CreatedAt   sql.TimeNullable `json:"created_at"`
	UpdatedAt   sql.TimeNullable `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at"  gorm:"index"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Policies    []Policy         `json:"policies"    gorm:"foreignKey:RoleID"`
}
