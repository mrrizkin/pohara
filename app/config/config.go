package config

import (
	"time"
)

type Config struct {
	App struct {
		Name    string `env:"APP_NAME,required"`
		Env     string `env:"APP_ENV,required"`
		URL     string `env:"APP_URL,required"`
		Port    int    `env:"APP_PORT,required"`
		Prefork bool   `env:"APP_PREFORK,default=false"`
		Debug   bool   `env:"APP_DEBUG,default=false"`
	}

	Site Site

	CSRF struct {
		Key        string `env:"CSRF_KEY,default=X-CSRF-Token"`
		CookieName string `env:"CSRF_COOKIE_NAME,default=fiber_csrf_token"`
		SameSite   string `env:"CSRF_SAME_SITE,default=Lax"`
		Secure     bool   `env:"CSRF_SECURE,default=false"`
		Session    bool   `env:"CSRF_SESSION,default=true"`
		HttpOnly   bool   `env:"CSRF_HTTP_ONLY,default=true"`
		Expiration int    `env:"CSRF_EXPIRATION,default=3600"`
	}

	View struct {
		Directory string `env:"VIEW_DIRECTORY,default=/views"`
		Extension string `env:"VIEW_EXTENSION,default=.html"`
		Cache     bool   `env:"VIEW_CACHE,default=true"`
	}

	StoragePath string `env:"STORAGE_PATH,default=storage"`
	SwaggerPath string `env:"SWAGGER_PATH,default=/docs/swagger.json"`
	CacheTTL    int    `env:"CACHE_TTL,default=300"`
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production" || c.App.Env == "prod"
}

func (c *Config) IsCacheView() bool {
	return c.View.Cache && c.IsProduction()
}

func (c *Config) CacheTTLSecond() time.Duration {
	return time.Duration(c.CacheTTL) * time.Second
}
