package cli

import (
	"github.com/mrrizkin/pohara/modules/database/migration"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type MigrateCmd struct {
	migrator *migration.Migrator
}

type MigrateCmdDeps struct {
	fx.In

	Migrator *migration.Migrator
}

func NewMigratorCmd(deps MigrateCmdDeps) *MigrateCmd {
	return &MigrateCmd{migrator: deps.Migrator}
}

func (*MigrateCmd) Name() string {
	return "migrate"
}

func (*MigrateCmd) Description() string {
	return "migrate database"
}

func (c *MigrateCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		if err := c.migrator.Migrate(); err != nil {
			cmd.PrintErrf("error running migration: %v\n", err)
			return
		}

		cmd.Println("Migration completed successfully")
		return
	}
	switch args[0] {
	case "rollback":
		if err := c.migrator.RollbackLastBatch(); err != nil {
			cmd.PrintErrf("Error rolling back migration: %v\n", err)
			return
		}
		cmd.Println("Rollback completed successfully")

	case "reset":
		if err := c.migrator.RollbackAll(); err != nil {
			cmd.PrintErrf("Error resetting migrations: %v\n", err)
			return
		}
		cmd.Println("Reset completed successfully")

	case "status":
		if err := c.migrator.Status(); err != nil {
			cmd.PrintErrf("Error getting migration status: %v\n", err)
			return
		}

	case "create":
		if len(args) < 2 {
			cmd.PrintErrf("Missing migration name\n")
			return
		}

		if err := c.migrator.CreateMigration(args[1]); err != nil {
			cmd.PrintErrf("Error creating migration: %v\n", err)
			return
		}
		cmd.Println("Migration created successfully")
	default:
		cmd.PrintErrf("Unknown command: %s\n", args[0])
		cmd.Println("Available commands: rollback, reset, status")
	}
}
