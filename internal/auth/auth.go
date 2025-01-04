package auth

import (
	"github.com/mrrizkin/pohara/internal/auth/delivery"
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/infrastructure/database"
	"github.com/mrrizkin/pohara/internal/server"
	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(
		delivery.Handler,
		server.AsWebRouter(delivery.WebRouter),
		server.AsApiRouter(delivery.ApiRouter),

		database.AsGormMigration(&entity.Role{}),
		database.AsGormMigration(&entity.Policy{}),
	),
)
