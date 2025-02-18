package database

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/common/config"
	"github.com/mrrizkin/pohara/modules/database/cli"
	dbConfig "github.com/mrrizkin/pohara/modules/database/config"
	"github.com/mrrizkin/pohara/modules/database/db"
	"github.com/mrrizkin/pohara/modules/database/migration"
	"github.com/mrrizkin/pohara/modules/database/repository"
)

var Module = fx.Module("database",
	config.Load(&dbConfig.Config{}),
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
