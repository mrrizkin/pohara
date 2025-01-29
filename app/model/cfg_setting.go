package model

type CfgSetting struct {
	ID uint `json:"id" db:"id"`
}

func (CfgSetting) TableName() string {
	return "cfg_setting"
}
