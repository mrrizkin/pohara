package cli

import (
	"github.com/mrrizkin/pohara/modules/logger"
	"github.com/spf13/cobra"
)

type Command interface {
	Name() string
	Description() string
	Run(cmd *cobra.Command, args []string)
}

type CLI struct {
	command *cobra.Command
	log     *logger.Logger
}

func New(logger *logger.Logger) *CLI {
	command := &cobra.Command{
		Use:   "Pohara",
		Short: "Pohara CLI",
	}
	return &CLI{
		command: command,
		log:     logger,
	}
}

func (c *CLI) addCommand(command Command) {
	c.command.AddCommand(&cobra.Command{
		Use:   command.Name(),
		Short: command.Description(),
		Run:   command.Run,
	})
}
