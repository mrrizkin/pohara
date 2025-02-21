package config

type Config struct {
	Title         string `env:"SITE_TITLE"`
	BaseURL       string `env:"SITE_BASE_URL"`
	BasePath      string `env:"SITE_BASE_PATH"`
	TrailingSlash bool   `env:"SITE_TRAILING_SLASH"`
	Favicon       string `env:"SITE_FAVICON"`
	Logo          string `env:"SITE_LOGO"`
	LogoWidth     int    `env:"SITE_LOGO_WIDTH"`
	LogoHeight    int    `env:"SITE_LOGO_HEIGHT"`
	LogoText      string `env:"SITE_LOGO_TEXT"`

	Settings struct {
		Pagination    int `env:"SITE_SETTINGS_PAGINATION"`
		SummaryLength int `env:"SITE_SETTINGS_SUMMARY_LENGTH"`
	}

	NavButton struct {
		Enable bool   `env:"SITE_NAV_BUTTON_ENABLE"`
		Label  string `env:"SITE_NAV_BUTTON_LABEL"`
		Link   string `env:"SITE_NAV_BUTTON_LINK"`
	}

	Params struct {
		ContactFormAction string `env:"SITE_PARAMS_CONTACT_FORM_ACTION"`
		TagManagerID      string `env:"SITE_PARAMS_TAG_MANAGER_ID"`
		FooterContent     string `env:"SITE_PARAMS_FOOTER_CONTENT"`
		Copyright         string `env:"SITE_PARAMS_COPYRIGHT"`
	}

	Metadata struct {
		MetaAuthor      string `env:"SITE_METADATA_META_AUTHOR"`
		MetaImage       string `env:"SITE_METADATA_META_IMAGE"`
		MetaDescription string `env:"SITE_METADATA_META_DESCRIPTION"`
	}

	StoragePath string `env:"STORAGE_PATH,default=storage"`
}
