package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
)

type CreateCfgSettingTable struct{}

func (m *CreateCfgSettingTable) ID() string {
	return "20250118143703_create_cfg_setting_table"
}

func (m *CreateCfgSettingTable) Up(schema *migration.Schema) {
	schema.Create("cfg_setting", func(table *migration.Blueprint) {
		table.ID()
		table.Text("site_name")
		table.Text("logo").Nullable()
		table.Timestamps()
	})
}

func (m *CreateCfgSettingTable) Down(schema *migration.Schema) {
	schema.Drop("cfg_setting")
}
