package service

import "go.uber.org/fx"

var Provide = fx.Provide(NewAuthService)
