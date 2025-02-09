package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
)

type CreateMUserSettingTable struct{}

func (m *CreateMUserSettingTable) ID() string {
	return "20250209165817_create_m_user_setting_table"
}

func (m *CreateMUserSettingTable) Up(schema *migration.Schema) {
	schema.Create("m_user_setting", func(table *migration.Blueprint) {
		table.ID()
		table.BigInteger("user_id")
		table.Text("language")
		table.Text("theme")
		table.Timestamps()
	})
}

func (m *CreateMUserSettingTable) Down(schema *migration.Schema) {
	schema.Drop("m_user_setting")
}
