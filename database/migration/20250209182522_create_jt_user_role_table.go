package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
)

type CreateJtUserRoleTable struct{}

func (m *CreateJtUserRoleTable) ID() string {
	return "20250209182522_create_jt_user_role_table"
}

func (m *CreateJtUserRoleTable) Up(schema *migration.Schema) {
	schema.CreateNotExist("jt_user_role", func(table *migration.Blueprint) {
		table.BigInteger("user_id").Index("idx_user_id").Foreign("m_user", "id").OnDelete("CASCADE")
		table.BigInteger("role_id").Index("idx_role_id").Foreign("m_role", "id").OnDelete("CASCADE")
		table.Primary("user_id", "role_id")
	})
}

func (m *CreateJtUserRoleTable) Down(schema *migration.Schema) {
	schema.DropExist("jt_user_role")
}
