package repository

import "go.uber.org/fx"

var Provide = fx.Provide(
	NewAuthRepository,
	NewUserRepository,
	NewSettingRepository,
	NewRoleRepository,
	NewPolicyRepository,
)
