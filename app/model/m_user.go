package model

import "github.com/mrrizkin/pohara/modules/common/sql"

type MUser struct {
	ID       uint               `json:"id"       db:"id"`
	Name     string             `json:"name"     db:"name"`
	Username string             `json:"username" db:"username"`
	Password string             `json:"-"        db:"password"`
	Email    sql.StringNullable `json:"email"    db:"email"`
}

func (MUser) TableName() string {
	return "m_user"
}
