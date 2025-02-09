package model

import "github.com/mrrizkin/pohara/modules/common/sql"

type MUserSetting struct {
	ID        uint             `json:"id"         gorm:"primaryKey"`
	UserID    uint             `json:"user_id"`
	Language  string           `json:"language"`
	Theme     string           `json:"theme"`
	CreatedAt sql.TimeNullable `json:"created_at"`
	UpdatedAt sql.TimeNullable `json:"updated_at"`
}

func (MUserSetting) TableName() string {
	return "m_user_setting"
}
