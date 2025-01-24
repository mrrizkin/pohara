package config

import (
	"go.uber.org/fx"
)

func Load(configs ...interface{}) fx.Option {
	for _, config := range configs {
		err := load(config)
		if err != nil {
			panic(err)
		}
	}

	return fx.Supply(configs...)
}
