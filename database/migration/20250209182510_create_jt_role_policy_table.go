package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
)

type CreateJtRolePolicyTable struct{}

func (m *CreateJtRolePolicyTable) ID() string {
	return "20250209182510_create_jt_role_policy_table"
}

func (m *CreateJtRolePolicyTable) Up(schema *migration.Schema) {
	schema.Create("jt_role_policy", func(table *migration.Blueprint) {
		table.ID()
		table.BigInteger("role_id")
		table.BigInteger("policy_id")
	})
}

func (m *CreateJtRolePolicyTable) Down(schema *migration.Schema) {
	schema.Drop("jt_role_policy")
}
