package migration

import "github.com/mrrizkin/pohara/modules/core/migration"

type CreateMRoleTable struct{}

func (m *CreateMRoleTable) ID() string {
	return "20250118103233_create_m_role_table"
}

func (m *CreateMRoleTable) Up(schema *migration.Schema) {
	schema.Create("m_role", func(table *migration.Blueprint) {
		table.ID()
		table.Text("name").Unique()
		table.Text("description")
		table.Timestamps()
	})
}

func (m *CreateMRoleTable) Down(schema *migration.Schema) {
	schema.Drop("m_role")
}
