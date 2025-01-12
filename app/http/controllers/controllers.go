package controllers

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/http/controllers/admin"
	"github.com/mrrizkin/pohara/app/http/controllers/api"
)

var Module = fx.Module("controllers",
	fx.Provide(
		NewWelcomeController,
		NewAuthController,

		api.NewUserController,

		admin.NewSettingController,
		admin.NewDashboardController,
	),
)
