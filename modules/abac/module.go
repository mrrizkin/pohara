package abac

import (
	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(NewAuthorizationService),
)
