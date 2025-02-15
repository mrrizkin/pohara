package migration

import "github.com/mrrizkin/pohara/modules/database/migration"

type CreateCfgPolicyTable struct{}

func (m *CreateCfgPolicyTable) ID() string {
	return "20250118103239_create_cfg_policy_table"
}

func (m *CreateCfgPolicyTable) Up(schema *migration.Schema) {
	schema.CreateNotExist("cfg_policy", func(table *migration.Blueprint) {
		table.ID()
		table.Text("name")
		table.Text("condition").Nullable()
		table.Text("action")
		table.Text("effect")
		table.Text("resource").Nullable()
		table.Timestamps()
	})
}

func (m *CreateCfgPolicyTable) Down(schema *migration.Schema) {
	schema.DropExist("cfg_policy")
}
