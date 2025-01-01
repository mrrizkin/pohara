package auth

import (
	"github.com/mrrizkin/pohara/internal/auth/delivery"
	"github.com/mrrizkin/pohara/internal/auth/entity"
	"github.com/mrrizkin/pohara/internal/infrastructure/database"
	"github.com/mrrizkin/pohara/internal/server"
	"github.com/mrrizkin/pohara/internal/web/template"
	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(
		delivery.Handler,
		server.AsWebRouter(delivery.WebRouter),
		server.AsApiRouter(delivery.ApiRouter),

		template.Extend(func() template.ExtendResult {
			return template.NewExtend(map[string]interface{}{
				"can": template.CustomIf("can", Can),
			})
		}),

		database.AsGormMigration(&entity.Role{}),
		database.AsGormMigration(&entity.Policy{}),
	),
)
