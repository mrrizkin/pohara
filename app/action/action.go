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
	PageUser = access.NewViewPageAction("page-user")

	PageSettingProfile      = access.NewViewPageAction("page-setting-profile")
	PageSettingAccount      = access.NewViewPageAction("page-setting-account")
	PageSettingAppearance   = access.NewViewPageAction("page-setting-appearance")
	PageSettingNotification = access.NewViewPageAction("page-setting-notification")
	PageSettingDisplay      = access.NewViewPageAction("page-setting-display")

	PageSystemSettingBranding     = access.NewViewPageAction("page-system-setting-branding")
	PageSystemSettingIntegration  = access.NewViewPageAction("page-system-setting-integration")
	PageSystemSettingGeneral      = access.NewViewPageAction("page-system-setting-general")
	PageSystemSettingSecurity     = access.NewViewPageAction("page-system-setting-security")
	PageSystemSettingLocalization = access.NewViewPageAction("page-system-setting-localization")
	PageSystemSettingNotification = access.NewViewPageAction("page-system-setting-notification")
	PageSystemSettingBackup       = access.NewViewPageAction("page-system-setting-backup")
	PageSystemSettingPurge        = access.NewViewPageAction("page-system-setting-purge")

	PageSystemAuthRole   = access.NewViewPageAction("page-system-auth-role")
	PageSystemAuthPolicy = access.NewViewPageAction("page-system-auth-policy")
)
