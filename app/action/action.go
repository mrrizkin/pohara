package action

import (
	"github.com/mrrizkin/pohara/modules/abac/access"
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
	PageUser = NewViewPageAction("page-user")

	PageSettingProfile      = NewViewPageAction("page-setting-profile")
	PageSettingAccount      = NewViewPageAction("page-setting-account")
	PageSettingAppearance   = NewViewPageAction("page-setting-appearance")
	PageSettingNotification = NewViewPageAction("page-setting-notification")
	PageSettingDisplay      = NewViewPageAction("page-setting-display")

	PageSystemSettingBranding     = NewViewPageAction("page-system-setting-branding")
	PageSystemSettingIntegration  = NewViewPageAction("page-system-setting-integration")
	PageSystemSettingGeneral      = NewViewPageAction("page-system-setting-general")
	PageSystemSettingSecurity     = NewViewPageAction("page-system-setting-security")
	PageSystemSettingLocalization = NewViewPageAction("page-system-setting-localization")
	PageSystemSettingNotification = NewViewPageAction("page-system-setting-notification")
	PageSystemSettingBackup       = NewViewPageAction("page-system-setting-backup")
	PageSystemSettingPurge        = NewViewPageAction("page-system-setting-purge")

	PageSystemAuthRole   = NewViewPageAction("page-system-auth-role")
	PageSystemAuthPolicy = NewViewPageAction("page-system-auth-policy")
)

func NewViewPageAction(resource string) access.Action {
	return access.NewResourceAction(resource, "view-page")
}
