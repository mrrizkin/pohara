package middleware

import "go.uber.org/fx"

var Provide = fx.Provide(
	NewAuthMiddleware,
	NewMenuMiddleware,
)
