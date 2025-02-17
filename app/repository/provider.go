package repository

import "go.uber.org/fx"

var Provide = fx.Provide(
	NewUserRepository,
	NewSettingRepository,
	NewRoleRepository,
)
