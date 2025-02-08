package cli

import (
	"os"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type Command interface {
	Name() string
	Description() string
	Run(cmd *cobra.Command, args []string)
}

type CLI struct {
	command *cobra.Command
	log     *logger.ZeroLog
}

type CommandLoaderDependencies struct {
	fx.In

	Cli      *CLI
	Commands []Command `group:"command"`
}

var Module = fx.Module("cli",
	fx.Provide(New),
	fx.Decorate(LoadCommand),
	fx.Invoke(Execute),
)

func New(config *config.Config, logger *logger.ZeroLog) *CLI {
	command := &cobra.Command{
		Use:   config.App.Name,
		Short: config.App.Name + " CLI",
	}
	command.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "running the server",
	})
	return &CLI{
		command: command,
		log:     logger,
	}
}

func LoadCommand(deps CommandLoaderDependencies) *CLI {
	for _, command := range deps.Commands {
		deps.Cli.AddCommand(command)
	}

	return deps.Cli
}

func AsCommand(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"command"`),
	)
}

func (c *CLI) AddCommand(command Command) {
	c.command.AddCommand(&cobra.Command{
		Use:   command.Name(),
		Short: command.Description(),
		Run:   command.Run,
	})
}

func Execute(c *CLI) {
	if err := c.command.Execute(); err != nil {
		c.log.Error("error execute command", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}
