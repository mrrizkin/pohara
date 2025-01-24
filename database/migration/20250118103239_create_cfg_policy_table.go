package migration

import "github.com/mrrizkin/pohara/modules/core/migration"

type CreateCfgPolicyTable struct{}

func (m *CreateCfgPolicyTable) ID() string {
	return "20250118103239_create_cfg_policy_table"
}

func (m *CreateCfgPolicyTable) Up(schema *migration.Schema) {
	schema.Create("cfg_policy", func(table *migration.Blueprint) {
		table.ID()
		table.Text("name")
		table.Json("conditions")
		table.Text("action")
		table.Text("effect")
		table.Text("resource")
		table.Timestamps()
	})
}

func (m *CreateCfgPolicyTable) Down(schema *migration.Schema) {
	schema.Drop("cfg_policy")
}
