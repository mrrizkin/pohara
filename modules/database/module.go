package database

import (
	"github.com/mrrizkin/pohara/modules/database/cli"
	"github.com/mrrizkin/pohara/modules/database/db"
	"github.com/mrrizkin/pohara/modules/database/migration"
	"github.com/mrrizkin/pohara/modules/database/repository"
	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(
		db.NewDatabase,
		repository.NewMigrationHistoryRepository,
		migration.NewMigrator,
		cli.NewMigratorCmd,
	),
	fx.Provide(
		cli.RegisterCommands,
	),
	fx.Decorate(migration.LoadMigration),
)
