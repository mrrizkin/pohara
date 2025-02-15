package migration

import "github.com/mrrizkin/pohara/modules/database/migration"

type CreateMUserAttributeTable struct{}

func (m *CreateMUserAttributeTable) ID() string {
	return "20250118103248_create_m_user_attribute_table"
}

func (m *CreateMUserAttributeTable) Up(schema *migration.Schema) {
	schema.CreateNotExist("m_user_attribute", func(table *migration.Blueprint) {
		table.ID()
		table.BigInteger("user_id").Foreign("m_user", "id").OnDelete("CASCADE")
		table.Text("location")
		table.Timestamps()
	})
}

func (m *CreateMUserAttributeTable) Down(schema *migration.Schema) {
	schema.DropExist("m_user_attribute")
}
