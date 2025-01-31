package migration

import (
	"github.com/mrrizkin/pohara/modules/core/migration"
)

type CreateCfgSettingTable struct{}

func (m *CreateCfgSettingTable) ID() string {
	return "20250118143703_create_cfg_setting_table"
}

func (m *CreateCfgSettingTable) Up(schema *migration.Schema) {
	schema.Create("cfg_setting", func(table *migration.Blueprint) {
		table.ID()
		table.Timestamps()
	})
}

func (m *CreateCfgSettingTable) Down(schema *migration.Schema) {
	schema.Drop("cfg_setting")
}
