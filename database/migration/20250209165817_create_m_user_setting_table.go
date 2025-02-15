package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
)

type CreateMUserSettingTable struct{}

func (m *CreateMUserSettingTable) ID() string {
	return "20250209165817_create_m_user_setting_table"
}

func (m *CreateMUserSettingTable) Up(schema *migration.Schema) {
	schema.CreateNotExist("m_user_setting", func(table *migration.Blueprint) {
		table.ID()
		table.BigInteger("user_id").Foreign("m_user", "id").OnDelete("CASCADE")
		table.Text("language")
		table.Text("theme")
		table.Timestamps()
	})
}

func (m *CreateMUserSettingTable) Down(schema *migration.Schema) {
	schema.DropExist("m_user_setting")
}
