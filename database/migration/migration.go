package migration

import (
	"github.com/mrrizkin/pohara/modules/core/migrator"
	"go.uber.org/fx"
)

var Module = fx.Module("migration",
	migrator.ProvideMigration(
		/** PLACEHOLDER **/
		&CreateMUserAttributeTable{},
		&CreateCfgPolicyTable{},
		&CreateMRoleTable{},
		&CreateMUserTable{},
	),
)
