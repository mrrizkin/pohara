package config

type Config struct {
	Name    string `env:"SERVER_NAME,required"`
	Env     string `env:"SERVER_ENV,required"`
	URL     string `env:"SERVER_URL,required"`
	Port    int    `env:"SERVER_PORT,required"`
	Prefork bool   `env:"SERVER_PREFORK,default=false"`
	Debug   bool   `env:"SERVER_DEBUG,default=false"`

	CSRF struct {
		Key        string `env:"CSRF_KEY,default=X-CSRF-Token"`
		CookieName string `env:"CSRF_COOKIE_NAME,default=fiber_csrf_token"`
		SameSite   string `env:"CSRF_SAME_SITE,default=Lax"`
		Secure     bool   `env:"CSRF_SECURE,default=false"`
		Session    bool   `env:"CSRF_SESSION,default=true"`
		HttpOnly   bool   `env:"CSRF_HTTP_ONLY,default=true"`
		Expiration int    `env:"CSRF_EXPIRATION,default=3600"`
	}

	SwaggerPath string `env:"SWAGGER_PATH,default=/docs/swagger.json"`
}

func (c *Config) IsProduction() bool {
	return c.Env == "production" || c.Env == "prod"
}
