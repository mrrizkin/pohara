package cli

import (
	"go.uber.org/fx"
)

type CommandLoaderDeps struct {
	fx.In

	Cli      *CLI
	Commands []Command `group:"command"`
}

func loadCommand(deps CommandLoaderDeps) *CLI {
	for _, command := range deps.Commands {
		deps.Cli.addCommand(command)
	}

	return deps.Cli
}

func execute(c *CLI) {
	if err := c.command.Execute(); err != nil {
		c.log.Error("error execute command", "error", err)
	}
}
