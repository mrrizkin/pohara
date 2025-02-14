package model

import "github.com/mrrizkin/pohara/modules/common/sql"

type CfgSetting struct {
	ID        uint               `json:"id"         gorm:"primaryKey"`
	SiteName  string             `json:"site_name"`
	Logo      sql.StringNullable `json:"logo"`
	CreatedAt   sql.TimeNullable `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt   sql.TimeNullable `json:"updated_at"  gorm:"autoUpdateTime"`
}

func (CfgSetting) TableName() string {
	return "cfg_setting"
}
