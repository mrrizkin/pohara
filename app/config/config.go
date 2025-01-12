package config

import "github.com/mrrizkin/pohara/modules/common/config"

var Module = config.Load(
	&App{},
	&Database{},
	&Session{},
	&Inertia{},
)
