package cli

import (
	"github.com/mrrizkin/pohara/modules/cli"
	"go.uber.org/fx"
)

type RegisterCommandsDeps struct {
	fx.In

	StartServerCmd *StartServerCmd
}

type RegisterCommandsResult struct {
	fx.Out

	StartServerCmd cli.Command `group:"command"`
}

func RegisterCommands(deps RegisterCommandsDeps) RegisterCommandsResult {
	return RegisterCommandsResult{
		StartServerCmd: deps.StartServerCmd,
	}
}
