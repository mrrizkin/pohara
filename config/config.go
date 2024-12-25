package config

import (
	"go.uber.org/fx"
)

func New() fx.Option {
	configs := []interface{}{
		&App{},
		&Database{},
		&Session{},
	}

	for _, config := range configs {
		load(config)
	}

	return fx.Supply(configs...)
}
