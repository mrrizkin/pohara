package config

import (
	"go.uber.org/fx"
)

func Load(cfgs ...interface{}) fx.Option {
	for _, config := range cfgs {
		load(config)
	}

	return fx.Supply(cfgs...)
}
