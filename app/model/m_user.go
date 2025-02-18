package model

import "github.com/mrrizkin/pohara/modules/common/sql"

type MUser struct {
	ID        uint               `json:"id"         gorm:"primaryKey"`
	Name      string             `json:"name"`
	Username  string             `json:"username"   gorm:"unique"`
	Password  string             `json:"-"`
	Email     sql.StringNullable `json:"email"      gorm:"unique"`
	CreatedAt sql.TimeNullable   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt sql.TimeNullable   `json:"updated_at" gorm:"autoUpdateTime"`
}

func (MUser) TableName() string {
	return "m_user"
}

func (m MUser) ResourceName() sql.StringNullable {
	return sql.String(m.TableName())
}
