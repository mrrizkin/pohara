package migration

import "github.com/mrrizkin/pohara/modules/database/migration"

type CreateMUserAttributeTable struct{}

func (m *CreateMUserAttributeTable) ID() string {
	return "20250118103248_create_m_user_attribute_table"
}

func (m *CreateMUserAttributeTable) Up(schema *migration.Schema) {
	schema.Create("m_user_attribute", func(table *migration.Blueprint) {
		table.ID()
		table.BigInteger("user_id")
		table.Text("location")
		table.Text("language")
		table.Timestamps()
	})
}

func (m *CreateMUserAttributeTable) Down(schema *migration.Schema) {
	schema.Drop("m_user_attribute")
}
