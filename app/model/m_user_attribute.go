package model

import "github.com/mrrizkin/pohara/modules/common/sql"

type MUserAttribute struct {
	ID        uint             `json:"id"         gorm:"primaryKey"`
	UserID    uint             `json:"user_id"`
	Location  string           `json:"location"`
	CreatedAt sql.TimeNullable `json:"created_at"`
	UpdatedAt sql.TimeNullable `json:"updated_at"`
}

func (MUserAttribute) TableName() string {
	return "m_user_attribute"
}
