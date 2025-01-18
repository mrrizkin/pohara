package migration

import (
	"gorm.io/gorm"
)

type CreateMRoleTable struct{}

func (m *CreateMRoleTable) ID() string {
	return "20250118103233_create_m_role_table"
}

func (m *CreateMRoleTable) Up(tx *gorm.DB) error {
	return tx.Exec(`CREATE TABLE IF NOT EXISTS m_role (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT
	)`).Error
}

func (m *CreateMRoleTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS m_role`).Error
}
