package migration

import (
	"sort"

	"go.uber.org/fx"
)

type MigrationLoaderDeps struct {
	fx.In

	Migrator   *Migrator
	Migrations []Migration `group:"migration"`
}

func LoadMigration(deps MigrationLoaderDeps) *Migrator {
	sortedMigrations := make([]Migration, len(deps.Migrations))
	copy(sortedMigrations, deps.Migrations)

	sort.Slice(sortedMigrations, func(i, j int) bool {
		return sortedMigrations[i].ID() < sortedMigrations[j].ID()
	})

	for _, migration := range sortedMigrations {
		deps.Migrator.AddMigration(migration)
	}

	return deps.Migrator
}

func ProvideMigration(migrations ...Migration) fx.Option {
	var options []fx.Option

	for _, migration := range migrations {
		m := migration
		options = append(options, fx.Provide(fx.Annotate(
			func() Migration { return m },
			fx.ResultTags(`group:"migration"`),
		)))
	}

	return fx.Options(options...)
}
