package migration

import (
	"gorm.io/gorm"
)

type CreateMUserAttributeTable struct{}

func (m *CreateMUserAttributeTable) ID() string {
	return "20250118103248_create_m_user_attribute_table"
}

func (m *CreateMUserAttributeTable) Up(tx *gorm.DB) error {
	return tx.Exec(`CREATE TABLE IF NOT EXISTS m_user_attribute (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		location TEXT,
		language TEXT
	)`).Error
}

func (m *CreateMUserAttributeTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS m_user_attribute`).Error
}
