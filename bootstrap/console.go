package bootstrap

import (
	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/database/migration"
	"github.com/mrrizkin/pohara/modules/cli"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/core"
	"github.com/mrrizkin/pohara/modules/core/migrator"
	"go.uber.org/fx"
)

func Console() *fx.App {
	return fx.New(
		config.Module,
		core.Module,
		hash.Module,
		repository.Module,
		migration.Module,
		migrator.Module,
		cli.Module,

		fx.NopLogger,
	)
}
