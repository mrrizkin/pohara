package migration

import (
	"gorm.io/gorm"
)

type CreateMUserTable struct{}

func (m *CreateMUserTable) ID() string {
	return "20250118103226_create_m_user_table"
}

func (m *CreateMUserTable) Up(tx *gorm.DB) error {
	return tx.Exec(`CREATE TABLE IF NOT EXISTS m_user (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	)`).Error
}

func (m *CreateMUserTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS m_user`).Error
}
