package model

import "github.com/mrrizkin/pohara/modules/common/sql"

type MUser struct {
	ID        uint               `json:"id"         gorm:"primaryKey"`
	Name      string             `json:"name"`
	Username  string             `json:"username"   gorm:"unique"`
	Password  string             `json:"-"`
	Email     sql.StringNullable `json:"email"      gorm:"unique"`
	CreatedAt sql.TimeNullable   `json:"created_at"`
	UpdatedAt sql.TimeNullable   `json:"updated_at"`
}

func (MUser) TableName() string {
	return "m_user"
}
