package bootstrap

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/database/migration"
	"github.com/mrrizkin/pohara/modules/cli"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/core"
	dbMigration "github.com/mrrizkin/pohara/modules/database/migration"
)

func Console() *fx.App {
	return fx.New(
		config.Module,
		core.Module,
		hash.Module,
		repository.Module,
		dbMigration.Module,
		migration.Module,
		cli.Module,

		fx.NopLogger,
	)
}
