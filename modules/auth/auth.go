package auth

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/auth/repository"
	"github.com/mrrizkin/pohara/modules/auth/service"
)

var Module = fx.Module("auth",
	fx.Provide(
		fx.Private,
		repository.NewAuthRepository,
	),
	fx.Provide(service.NewAuthService),
)
