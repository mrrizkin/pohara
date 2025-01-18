package migration

import (
	"gorm.io/gorm"
)

type CreateCfgSettingTable struct{}

func (m *CreateCfgSettingTable) ID() string {
	return "20250118143703_create_cfg_setting_table"
}

func (m *CreateCfgSettingTable) Up(tx *gorm.DB) error {
	return tx.Exec(`CREATE TABLE IF NOT EXISTS cfg_setting (
		id BIGSERIAL PRIMARY KEY
	)`).Error
}

func (m *CreateCfgSettingTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS cfg_setting`).Error
}
