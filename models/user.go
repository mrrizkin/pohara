package models

import (
	"github.com/mrrizkin/pohara/module/database/sql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint               `json:"id"         gorm:"primary_key"`
	CreatedAt sql.TimeNullable   `json:"created_at"`
	UpdatedAt sql.TimeNullable   `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `json:"deleted_at" gorm:"index"`
	Username  sql.StringNullable `json:"username"   gorm:"unique;not null;index"`
	Password  sql.StringNullable `json:"password"`
	Name      sql.StringNullable `json:"name"`
	Email     sql.StringNullable `json:"email"`
}
