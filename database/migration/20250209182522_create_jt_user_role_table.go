package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
)

type CreateJtUserRoleTable struct{}

func (m *CreateJtUserRoleTable) ID() string {
	return "20250209182522_create_jt_user_role_table"
}

func (m *CreateJtUserRoleTable) Up(schema *migration.Schema) {
	schema.Create("jt_user_role", func(table *migration.Blueprint) {
		table.ID()
		table.BigInteger("user_id")
		table.BigInteger("role_id")
	})
}

func (m *CreateJtUserRoleTable) Down(schema *migration.Schema) {
	schema.Drop("jt_user_role")
}
