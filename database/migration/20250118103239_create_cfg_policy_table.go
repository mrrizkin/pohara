package migration

import (
	"gorm.io/gorm"
)

type CreateCfgPolicyTable struct{}

func (m *CreateCfgPolicyTable) ID() string {
	return "20250118103239_create_cfg_policy_table"
}

func (m *CreateCfgPolicyTable) Up(tx *gorm.DB) error {
	return tx.Exec(`CREATE TABLE IF NOT EXISTS cfg_policy (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		conditions JSONB NOT NULL,
		action TEXT NOT NULL,
		effect TEXT NOT NULL,
		resource TEXT NOT NULL
	)`).Error
}

func (m *CreateCfgPolicyTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS cfg_policy`).Error
}
