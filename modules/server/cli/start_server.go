package cli

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/logger"
	"github.com/mrrizkin/pohara/modules/server/config"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type StartServerCmd struct {
	app    *fiber.App
	config *config.Config
	log    *logger.Logger
	lc     fx.Lifecycle
}

type StartServerCmdDeps struct {
	fx.In

	App    *fiber.App
	Config *config.Config
	Logger *logger.Logger
}

func NewStartServerCmd(lc fx.Lifecycle, deps StartServerCmdDeps) *StartServerCmd {
	return &StartServerCmd{
		app:    deps.App,
		config: deps.Config,
		log:    deps.Logger,
		lc:     lc,
	}
}

func (*StartServerCmd) Name() string {
	return "serve"
}

func (*StartServerCmd) Description() string {
	return "start server"
}

func (c *StartServerCmd) Run(cmd *cobra.Command, args []string) {
	c.lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := c.app.Listen(fmt.Sprintf(":%d", c.config.Port)); err != nil {
					c.log.Fatal("failed to start server", "error", err)
					ctx.Done()
				}
			}()
			c.log.Info("server started", "app_name", c.config.Name, "port", c.config.Port)
			return nil
		},
		OnStop: func(context.Context) error {
			return c.app.Shutdown()
		},
	})
}
