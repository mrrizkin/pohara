package cli

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/cli"
)

type RegisterCommandsDeps struct {
	fx.In

	MigrateCmd *MigrateCmd
}

type RegisterCommandsResult struct {
	fx.Out

	MigrateCmd cli.Command `group:"command"`
}

func RegisterCommands(deps RegisterCommandsDeps) RegisterCommandsResult {
	return RegisterCommandsResult{
		MigrateCmd: deps.MigrateCmd,
	}
}
