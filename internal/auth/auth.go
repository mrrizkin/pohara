package auth

import (
	"github.com/mrrizkin/pohara/internal/auth/delivery"
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/auth/service"
	"github.com/mrrizkin/pohara/internal/infrastructure/database"
	"github.com/mrrizkin/pohara/internal/server"
	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(
		fx.Private,
		service.NewPolicyDecision,
		service.NewPolicyEnforcement,
	),

	fx.Provide(
		delivery.Handler,
		server.AsWebRouter(delivery.WebRouter),
		server.AsApiRouter(delivery.ApiRouter),

		service.NewService,

		database.AsGormMigration(&entity.Role{}, &entity.Policy{}),
	),
)
