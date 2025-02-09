package migration

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
	"go.uber.org/fx"
)

var Module = fx.Module("migration",
	migration.ProvideMigration(
		&CreateMUserTable{},
		&CreateMRoleTable{},
		&CreateCfgPolicyTable{},
		&CreateMUserAttributeTable{},
		&CreateCfgSettingTable{},
		&CreateMUserSettingTable{},
		&CreateJtRolePolicyTable{},
		&CreateJtUserRoleTable{},
		/** PLACEHOLDER **/
	),
)
