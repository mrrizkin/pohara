package entity

import (
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/common/sql"
)

type User struct {
	entity.BaseSubject
	Username sql.StringNullable `json:"username"   gorm:"unique;not null;index"`
	Password sql.StringNullable `json:"-"`
	Email    sql.StringNullable `json:"email"`
}
