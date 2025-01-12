package model

import "github.com/mrrizkin/pohara/modules/common/sql"

type User struct {
	Name     string             `json:"name"     gorm:"not null"`
	Username string             `json:"username" gorm:"unique;not null;index"`
	Password string             `json:"-"        gorm:"not null"`
	Email    sql.StringNullable `json:"email"    gorm:"not null"`
}
