package controllers

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/http/controllers/admin"
	"github.com/mrrizkin/pohara/app/http/controllers/admin/setting"
	"github.com/mrrizkin/pohara/app/http/controllers/admin/setting/system"
	"github.com/mrrizkin/pohara/app/http/controllers/api"
)

var Module = fx.Module("controllers",
	fx.Provide(
		NewClientPageController,
		NewAuthController,

		NewSetupController,

		api.NewUserController,

		admin.NewDashboardController,
		admin.NewUserController,

		setting.NewSettingController,

		system.NewSystemSettingController,
		system.NewAuthController,
	),
)
