package config

import "go.uber.org/fx"

func Load(cfgs ...interface{}) fx.Option {
	for _, config := range cfgs {
		err := load(config)
		if err != nil {
			panic(err)
		}
	}

	return fx.Supply(cfgs...)
}
