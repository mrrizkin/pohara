package migration

import (
	"github.com/mrrizkin/pohara/modules/cli"
	"github.com/spf13/cobra"
)

type MigratorCommand struct {
	migrator *Migrator
}

func NewMigratorCommand(migrator *Migrator) cli.Command {
	return &MigratorCommand{
		migrator: migrator,
	}
}

func (m *MigratorCommand) Name() string {
	return "migrate"
}

func (m *MigratorCommand) Description() string {
	return "migrate database"
}

func (m *MigratorCommand) Run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		if err := m.migrator.migrate(); err != nil {
			cmd.PrintErrf("error running migration: %v\n", err)
			return
		}

		cmd.Println("Migration completed successfully")
		return
	}
	switch args[0] {
	case "rollback":
		if err := m.migrator.rollbackLastBatch(); err != nil {
			cmd.PrintErrf("Error rolling back migration: %v\n", err)
			return
		}
		cmd.Println("Rollback completed successfully")

	case "reset":
		if err := m.migrator.rollbackAll(); err != nil {
			cmd.PrintErrf("Error resetting migrations: %v\n", err)
			return
		}
		cmd.Println("Reset completed successfully")

	case "status":
		if err := m.migrator.status(); err != nil {
			cmd.PrintErrf("Error getting migration status: %v\n", err)
			return
		}

	case "create":
		if len(args) < 2 {
			cmd.PrintErrf("Missing migration name\n")
			return
		}

		if err := m.migrator.createMigration(args[1]); err != nil {
			cmd.PrintErrf("Error creating migration: %v\n", err)
			return
		}
		cmd.Println("Migration created successfully")
	default:
		cmd.PrintErrf("Unknown command: %s\n", args[0])
		cmd.Println("Available commands: rollback, reset, status")
	}
}
