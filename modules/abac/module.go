package abac

import (
	"github.com/mrrizkin/pohara/modules/abac/service"
	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(service.NewAuthorizationService),
)
