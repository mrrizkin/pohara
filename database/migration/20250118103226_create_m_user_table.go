package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
)

type CreateMUserTable struct{}

func (m *CreateMUserTable) ID() string {
	return "20250118103226_create_m_user_table"
}

func (m *CreateMUserTable) Up(schema *migration.Schema) {
	schema.Create("m_user", func(table *migration.Blueprint) {
		table.ID()
		table.Text("name")
		table.Text("username").Unique()
		table.Text("password")
		table.Text("email").Unique()
		table.Timestamps()
	})
}

func (m *CreateMUserTable) Down(schema *migration.Schema) {
	schema.Drop("m_user")
}
