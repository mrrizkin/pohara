package action

import (
	"github.com/mrrizkin/pohara/modules/auth/access"
)

var (
	SpecialAll = access.NewResourceAction("special", "all")

	All    = access.NewResourceAction("common", "all")
	Create = access.NewResourceAction("common", "create")
	Read   = access.NewResourceAction("common", "read")
	View   = access.NewResourceAction("common", "view")
	Update = access.NewResourceAction("common", "update")
	Delete = access.NewResourceAction("common", "delete")

	/* you can integrate the page access to inertia in app/htt/middlware/menu */
	PageIntegration = access.NewViewPageAction("page-integration")

	PageSettingProfile      = access.NewViewPageAction("page-setting-profile")
	PageSettingAccount      = access.NewViewPageAction("page-setting-account")
	PageSettingAppearance   = access.NewViewPageAction("page-setting-appearance")
	PageSettingNotification = access.NewViewPageAction("page-setting-notification")
	PageSettingDisplay      = access.NewViewPageAction("page-setting-display")

	PageUser = access.NewViewPageAction("page-user")
)
