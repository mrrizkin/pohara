package model

type CfgSetting struct {
	ID uint `json:"id"         gorm:"primaryKey"`
}

func (CfgSetting) TableName() string {
	return "cfg_setting"
}
