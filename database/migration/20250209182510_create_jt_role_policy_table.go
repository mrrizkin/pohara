package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
)

type CreateJtRolePolicyTable struct{}

func (m *CreateJtRolePolicyTable) ID() string {
	return "20250209182510_create_jt_role_policy_table"
}

func (m *CreateJtRolePolicyTable) Up(schema *migration.Schema) {
	schema.CreateNotExist("jt_role_policy", func(table *migration.Blueprint) {
		table.BigInteger("role_id").Index("idx_role_id").Foreign("m_role", "id").OnDelete("CASCADE")
		table.BigInteger("policy_id").Index("idx_policy_id").Foreign("cfg_policy", "id").OnDelete("CASCADE")
		table.Primary("role_id", "policy_id")
	})
}

func (m *CreateJtRolePolicyTable) Down(schema *migration.Schema) {
	schema.DropExist("jt_role_policy")
}
